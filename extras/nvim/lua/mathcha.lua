local M = {}
---@type table<integer, State>
M._states = {}

---@alias NodeID string

---@class State
---@field buf integer
---@field query vim.treesitter.Query
---@field pending TSNode[]
---@field running vim.SystemObj?
local State = {}

local ns_id = vim.api.nvim_create_namespace("mathcha")

---zero indexed
---@generic T
---@param iter fun(any, any):T
---@return fun(): integer, T
local function enumerate(iter)
	local i = -1 -- the things I do for zero-indexing
	return function()
		local value = iter()
		if value == nil then return nil end
		i = i + 1
		return i, value
	end
end

---@param cmd string[]
---@param opt table
---@param win_opt vim.api.keyset.win_config
---@return integer?
---@return string?
local function jobstart_in_floating_win(cmd, opt, win_opt)
	assert(opt.term, "jobstart_in_floating_win called with no term flag")
	local buf = vim.api.nvim_create_buf(false, true)
	if buf == 0 then
		return nil, "could not create buffer"
	end

	local win = vim.api.nvim_open_win(buf, true, win_opt)
	if win == 0 then
		return nil, "could not create floating win"
	end
	vim.bo[buf].bufhidden = "wipe"
	vim.wo[win][0].winfixbuf = true

	return vim.fn.jobstart(cmd, opt), nil
end

---@param bufnr integer
---@return State? state
---@return string? error
function State.new(bufnr)
	vim.api.nvim_buf_clear_namespace(bufnr, ns_id, 0, -1)
	local state = setmetatable({
		buf = bufnr,
		jobs = {},
		pending = {},
	}, { __index = State })

	local err
	-- NOTE: we match the whole block including `$$` because otherwise we get
	-- each line separately for some reason.
	--
	-- It is also slightly annoying to match children of the latex_block, since
	-- those belong to the injected latex tree (and I don't want to juggle that
	-- many langs if I can avoid it), so manually skipping `$$` is the easiest
	-- way rn
	state.query, err = vim.treesitter.query.parse('markdown_inline', [[ (latex_block) @latex ]])
	state.equation_query = vim.treesitter.query.parse('markdown_inline', [[]])
	local md_inline_tree = vim.treesitter.get_parser(bufnr, "markdown")
		:children()["markdown_inline"]

	if md_inline_tree == nil then
		return nil, err
	end

	md_inline_tree:register_cbs({
		on_changedtree = function(_, tree)
			vim.schedule(function()
				-- FIXME: this forces a reload on all renders, I need to
				-- listen to on_bytes or nvim_buf_attach
				state:_parse_and_render({ [1] = tree })
			end)
		end,
		on_detach = function()
			vim.api.nvim_buf_clear_namespace(bufnr, ns_id, 0, -1)
			M._states[bufnr] = nil
		end
	}, true)

	md_inline_tree:parse(false, function(e, trees)
		if err ~= nil or trees == nil then
			err = ("failed to parse markdown_inline: %s"):format(e)
		else
			state:_parse_and_render(trees)
		end
	end)

	if err ~= nil then
		return nil, err
	end

	return state
end

---@param trees table<integer, TSTree>
function State:_parse_and_render(trees)
	for _, tree in pairs(trees) do
		for _, match, _ in self.query:iter_matches(tree:root(), self.buf, 0, -1) do
			for _, nodes in pairs(match) do
				vim.list_extend(self.pending, nodes)
			end
		end
	end
	self:update_conceal()
end

function State:update_conceal()
	local node = table.remove(self.pending, 1)
	while node and node:has_changes() do
		node = table.remove(self.pending, 1)
	end

	if node == nil then
		return
	end

	local start_row, _, end_row, _ = node:range()
	start_row = start_row + 1

	self.running = vim.system({ 'mathcha', '-render' }, {
		stdin = vim.api.nvim_buf_get_lines(self.buf, start_row, end_row, false)
	}, function(obj)
		self.running = nil
		vim.schedule(function()
			if obj.code ~= 0 then
				vim.notify("mathcha failed: " .. obj.stderr, vim.log.levels.ERROR)
			end

			local lines = enumerate(vim.gsplit(obj.stdout, '\n', { plain = true }))

			local virt_lines = {}
			for _, line in lines do
				table.insert(virt_lines, { { line } })
			end

			vim.api.nvim_buf_clear_namespace(self.buf, ns_id, start_row, end_row + 1)

			vim.api.nvim_buf_set_extmark(self.buf, ns_id, start_row, 0, {
				invalidate = true,
				conceal_lines = "",
				end_row = end_row - 1,
			})
			vim.api.nvim_buf_set_extmark(self.buf, ns_id, end_row, 0, {
				invalidate = true,
				virt_lines = virt_lines,
				virt_lines_above = true,
			})

			self:update_conceal()
		end)
	end)
end

function State:reset_marks()
	vim.api.nvim_buf_clear_namespace(self.buf, ns_id, 0, -1)
end

local function win_size()
	local ui = vim.api.nvim_list_uis()[1]
	local width = math.floor(ui.width * 0.9)
	local height = math.floor(ui.height * 0.9)
	return {
		width = width,
		height = height,
		col = math.floor((ui.width - width) / 2),
		row = math.floor((ui.height - height) / 2),
	}
end

-- see |on_stdout|
local function is_stdout_EOF_marker(l)
	return #l == 1 and l[1] == ""
end

---@return boolean
function State:open_editor()
	local cursor = vim.api.nvim_win_get_cursor(0)
	-- TODO: async?
	vim.treesitter.get_parser(self.buf):parse({ cursor[1] - 1, cursor[1] })

	local md_inline_tree = vim.treesitter.get_parser(self.buf, "markdown")
		:children()["markdown_inline"]

	local cur_range = { cursor[1] - 1, cursor[2], cursor[1] - 1, cursor[2] + 1 }
	local latex_node = md_inline_tree:named_node_for_range(cur_range, {})

	-- Note that "latex_block" is a leaf node, if this ever changes upstream
	-- this will break
	if not latex_node or latex_node:type() ~= "latex_block" then
		vim.notify("No node at cursor", vim.log.levels.WARN)
		return false
	end

	local start_row, start_col, end_row, _ = latex_node:range()
	start_row = start_row + 1 -- skip "$$"
	local latex_text = vim.api.nvim_buf_get_lines(self.buf, start_row, end_row, true)

	-- HACK: sync IO bad; and I couldn't get both TUI and passing via stdin to work
	-- at the same time
	local in_path = vim.fn.tempname()
	local fh, err = io.open(in_path, 'w')
	if err then
		vim.notify("could not open temp file: " .. err, vim.log.levels.ERROR)
		return false
	end
	assert(fh)
	_, err = fh:write(unpack(latex_text))
	if err then
		vim.notify("could not write to temp file: " .. err, vim.log.levels.ERROR)
		return false
	end
	fh:close()

	local sizes = win_size()
	local win_opt = {
		relative = "editor",
		width = sizes.width,
		height = sizes.height,
		col = sizes.col,
		row = sizes.row,
		style = "minimal",
		border = "rounded",
	}
	local out_marker_found = false
	local out_buf = {}
	local editor_cmd = { "mathcha", "-printout", "-f", in_path }
	local job
	job, err = jobstart_in_floating_win(editor_cmd, {
		term = true,
		on_exit = vim.schedule_wrap(function(_, code)
			-- TODO: sync bad
			os.remove(in_path)
			-- TODO is it better to use the TSNode for range + has_changes check
			if code ~= 0 then
				vim.notify("mathcha returned non-zero exit code " .. tostring(code), vim.log.levels.ERROR)
			elseif next(out_buf) then
				vim.api.nvim_buf_set_text(self.buf, start_row, start_col, end_row - 1, -1, out_buf)
			end
			vim.fn.chanclose(assert(job))
		end),
		on_stdout = function(_, data)
			---@cast data string[]
			if out_marker_found and not is_stdout_EOF_marker(data) then
				vim.list_extend(out_buf, data)
				return
			end
			for i, line in ipairs(data) do
				-- cursed magic string
				local MAGIC = '!mAtHcHa!'
				local mark = string.sub(line, 1, #MAGIC)
				-- HACK: I don't even know if this is consistent behavior
				if mark == MAGIC then
					local latex = string.sub(line, #MAGIC + 1)
					out_marker_found = true
					table.insert(out_buf, latex)
					vim.list_extend(out_buf, data, i + 1)
					return
				end
			end
		end,
	}, win_opt)

	if not job or err then
		vim.notify(assert(err), vim.log.levels.ERROR)
	end

	return true
end

function M.attach(bufnr)
	local buf = bufnr or vim.fn.bufnr()
	if buf == -1 then
		error("invalid bufnr " .. tostring(bufnr))
	end
	-- TODO: cleanup state on buf delete
	if not M._states[buf] then
		local err
		M._states[buf], err = State.new(buf)
		if err then
			vim.notify_once(err, vim.log.levels.ERROR)
		end
	end

	vim.keymap.set("n", "<localleader>i", M.open_editor, { buffer = buf })
end

function M.open_editor()
	local state = M._states[vim.fn.bufnr()]
	if not state then
		vim.notify("Is mathcha attached to this buffer?", vim.log.levels.WARN)
		return
	end
	state:open_editor()
end

function M.instance(bufnr)
	return M._states[vim.fn.bufnr(bufnr)]
end

-- for testing
function M.unload()
	for buf, _ in pairs(M._states) do
		vim.treesitter.stop(buf)
		vim.api.nvim_buf_clear_namespace(buf, ns_id, 0, -1)
	end
	M._states = {}
end

M.State = State

return M
