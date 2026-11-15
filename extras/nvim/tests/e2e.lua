-- End-to-end tests for the mathcha nvim plugin.
--
-- Run from extras/nvim with:
--
--   nvim --headless -u NONE -c 'luafile tests/e2e.lua'
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

---@param actual any
---@param expected any
---@param what string
local function assert_eq(actual, expected, what)
	if not vim.deep_equal(actual, expected) then
		error(string.format(
			"%s: expected %s, got %s",
			what,
			vim.inspect(expected),
			vim.inspect(actual)
		), 2)
	end
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

		vim.cmd.edit(vim.fn.fnameescape(fixture))
		local buf = vim.api.nvim_get_current_buf()
		assert_eq(vim.bo[buf].filetype, "markdown", "filetype")

		local ok, args = await_event("User", "MathchaJobsDone", 5000)
		if not ok then
			error(args, 0)
		end
		assert_eq(type(args), "table", "MathchaJobsDone args")

		local state = mathcha.instance(buf)
		if not state then
			error("mathcha did not attach to buffer " .. buf, 2)
		end
		assert_eq(state.buf, buf, "state.buf")
		assert_eq(#state.pending, 0, "state.pending")
		assert_eq(state.conceal_job, nil, "state.conceal_job")

		local marks = mathcha_extmarks(buf)
		assert_eq(#marks, 1, "extmark count")
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

local main = coroutine.create(function()
	-- tests must be run from the plugin root so the ftplugin is found
	local root = assert(vim.env.PWD, "PWD is not set")
	vim.opt.runtimepath:prepend(root)
	-- -u NONE means neither filetype detection nor ftplugins are enabled
	vim.cmd("filetype plugin on")

	for _, test in ipairs(tests) do
		log("test: %s", test.name)
		test.fn()
		log("ok: %s", test.name)
	end
end)

---@param ... any
resume_main = function(...)
	local ok, err = coroutine.resume(main, ...)
	if not ok then
		log("FAIL: %s", tostring(err))
		if headless then
			vim.cmd("cquit 1")
		end
		return
	end
	if coroutine.status(main) == "dead" then
		log("all tests passed")
		if headless then
			vim.cmd("qa!")
		end
	end
end

resume_main()
