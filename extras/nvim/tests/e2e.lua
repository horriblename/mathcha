-- End-to-end tests for the mathcha nvim plugin.
--
-- Run from extras/nvim with:
--
--   nvim --clean --headless -c 'luafile tests/e2e.lua'
--
-- The script adds PWD to 'runtimepath', so it must be run from the plugin
-- root. It quits nvim on success and exits with a non-zero code on failure.

local function log(fmt, ...)
	io.stdout:write(string.format(fmt, ...) .. "\n")
end
local headless = #vim.api.nvim_list_uis() == 0
if not headless then
	log = function(fmt, ...)
		vim.notify(string.format(fmt, ...))
	end
end

local results = {
	passed = 0,
	failed = 0,
	failures = {},
}
local current_test

-- Defined below, forward declared so async helpers can resume the top-level
-- coroutine and report failures through a single path.
local resume_main

---Wait for an autocmd from inside the main coroutine.
---
---The autocmd/timer resumes the main coroutine, so this must be called while
---the main coroutine is running.
---@return boolean ok, vim.api.keyset.create_autocmd.callback_args|string args
local function await_event(event, pattern, timeout_ms)
	timeout_ms = timeout_ms or 2000
	local co = assert(coroutine.running(), "await_event must run in a coroutine")
	local timer ---@type uv.uv_timer_t

	local au = vim.api.nvim_create_autocmd(event, {
		pattern = pattern,
		once = true,
		-- vim.schedule works around possible nested autocmd issues
		callback = vim.schedule_wrap(function(args)
			timer:stop()
			resume_main(true, args)
		end),
	})

	timer = vim.defer_fn(function()
		vim.api.nvim_del_autocmd(au)
		resume_main(false, "timed out waiting for " .. event .. " " .. pattern)
	end, timeout_ms)

	log("waiting for autocmd %s %s", event, pattern)
	return coroutine.yield()
end

---@param msg string
local function record_failure(msg)
	if current_test then
		msg = string.format("[%s] %s", current_test, msg)
	end
	results.failed = results.failed + 1
	results.failures[#results.failures + 1] = msg
	log("FAIL: %s", msg)
end

local function summarize()
	log("")
	log("summary: %d passed, %d failed", results.passed, results.failed)
	for _, msg in ipairs(results.failures) do
		log("  failed: %s", msg)
	end
end

---Record an assertion result instead of aborting the test run.
---@param actual any
---@param expected any
---@param what string
---@return boolean passed
local function assert_eq(actual, expected, what)
	if vim.deep_equal(actual, expected) then
		results.passed = results.passed + 1
		return true
	end
	record_failure(string.format(
		"%s: expected %s, got %s",
		what,
		vim.inspect(expected),
		vim.inspect(actual)
	))
	return false
end

---@param buf integer
---@return integer[][] marks with details
local function mathcha_extmarks(buf)
	local ns = vim.api.nvim_get_namespaces()["mathcha"]
	if not ns then
		return {}
	end
	return vim.api.nvim_buf_get_extmarks(buf, ns, 0, -1, { details = true })
end

---Plain-text lines from a virt_lines extmark detail.
---@param virt_lines table
---@return string[][]
local function virt_lines_text(virt_lines)
	local lines = {}
	for i, virt_line in ipairs(virt_lines) do
		lines[i] = {}
		for _, chunk in ipairs(virt_line) do
			lines[i][#lines[i] + 1] = chunk[1]
		end
	end
	return lines
end

---Open a fixture file, reusing the buffer if it is already open and resetting
---it to the file's on-disk contents.
---
---Reusing a single markdown buffer avoids juggling multiple markdown
---highlights/parsers between tests.
---@param path string
---@return integer buf
---@return boolean changed
local function open_fixture(path)
	local existing = vim.fn.bufnr(path)
	if existing == -1 then
		vim.cmd.edit(vim.fn.fnameescape(path))
		return vim.api.nvim_get_current_buf(), false
	end

	vim.api.nvim_set_current_buf(existing)
	local on_disk = vim.fn.readfile(path)
	if vim.deep_equal(vim.api.nvim_buf_get_lines(existing, 0, -1, false), on_disk) then
		return existing, false
	end
	vim.api.nvim_buf_set_lines(existing, 0, -1, false, on_disk)
	return existing, true
end

---@return boolean ok
local function await_jobs_done(timeout_ms)
	local ok, args = await_event("User", "MathchaJobsDone", timeout_ms or 5000)
	if not ok then
		record_failure(args)
		return false
	end
	return assert_eq(type(args), "table", "MathchaJobsDone args")
end

local tests = {}

tests[#tests + 1] = {
	name = "renders a $$ block into virtual lines",
	fn = function()
		local mathcha = require("mathcha")
		local fixture = vim.fs.joinpath(
			assert(vim.env.PWD, "PWD is not set"),
			"tests",
			"fixtures",
			"basic.md"
		)

		local buf = open_fixture(fixture)
		assert_eq(vim.bo[buf].filetype, "markdown", "filetype")

		if not await_jobs_done() then
			return
		end

		local state = mathcha.instance(buf)
		if not state then
			record_failure("mathcha did not attach to buffer " .. buf)
			return
		end
		assert_eq(state.buf, buf, "state.buf")
		assert_eq(#state.pending, 0, "state.pending")
		assert_eq(state.conceal_job, nil, "state.conceal_job")

		local marks = mathcha_extmarks(buf)
		if not assert_eq(#marks, 1, "extmark count") then
			return
		end
		local details = marks[1][4]
		assert_eq(details.invalidate, true, "extmark.invalidate")
		assert_eq(details.virt_lines_above, true, "extmark.virt_lines_above")
		assert_eq(
			virt_lines_text(details.virt_lines),
			{ { "" }, { "x² + y²" } },
			"rendered virtual lines"
		)
		assert_eq(details.virt_lines[1][1][2], "Special", "marker highlight")
		assert_eq(details.virt_lines[2][1][2], "TabLineSel", "render highlight")
	end,
}

tests[#tests + 1] = {
	name = "rerenders an existing $$ block after it is edited",
	fn = function()
		local mathcha = require("mathcha")
		local fixture = vim.fs.joinpath(
			assert(vim.env.PWD, "PWD is not set"),
			"tests",
			"fixtures",
			"basic.md"
		)

		local buf, changed = open_fixture(fixture)
		assert_eq(vim.bo[buf].filetype, "markdown", "filetype")

		-- if the buffer was already rendered by another test, there is no new
		-- MathchaJobsDone event to wait for
		if changed or not mathcha.instance(buf) or #mathcha_extmarks(buf) == 0 then
			if not await_jobs_done() then
				return
			end
		end

		local marks = mathcha_extmarks(buf)
		if not assert_eq(#marks, 1, "initial extmark count") then
			return
		end
		assert_eq(
			virt_lines_text(marks[1][4].virt_lines),
			{ { "" }, { "x² + y²" } },
			"initial rendered virtual lines"
		)

		-- replace the body of the existing latex block
		vim.api.nvim_buf_set_lines(buf, 3, 4, false, { "a^2 + b^2" })

		if not await_jobs_done() then
			return
		end

		local state = mathcha.instance(buf)
		if not state then
			record_failure("mathcha did not attach to buffer " .. buf)
			return
		end
		assert_eq(#state.pending, 0, "state.pending")
		assert_eq(state.conceal_job, nil, "state.conceal_job")

		marks = mathcha_extmarks(buf)
		if not assert_eq(#marks, 1, "extmark count after edit") then
			return
		end
		assert_eq(
			virt_lines_text(marks[1][4].virt_lines),
			{ { "" }, { "a² + b²" } },
			"rerendered virtual lines"
		)
	end,
}

local main = coroutine.create(function()
	-- tests must be run from the plugin root so the ftplugin is found
	local root = assert(vim.env.PWD, "PWD is not set")
	vim.opt.runtimepath:prepend(root)
	-- --clean skips user config, so filetype plugins are not enabled by default
	vim.cmd("filetype plugin on")

	for _, test in ipairs(tests) do
		current_test = test.name
		local failures_before = results.failed
		log("test: %s", test.name)
		test.fn()
		if results.failed == failures_before then
			log("ok: %s", test.name)
		else
			log("FAIL: %s", test.name)
		end
		current_test = nil
	end
end)

---@param ... any
resume_main = function(...)
	local ok, err = coroutine.resume(main, ...)
	if not ok then
		record_failure("unexpected error: " .. tostring(err))
		summarize()
		if headless then
			vim.cmd("cquit 1")
		end
		return
	end
	if coroutine.status(main) == "dead" then
		summarize()
		if results.failed == 0 then
			log("all tests passed")
			if headless then
				vim.cmd("qa!")
			end
		elseif headless then
			vim.cmd("cquit 1")
		end
	end
end

resume_main()
