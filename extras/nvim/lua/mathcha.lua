local M = {}

---@class State
---@field buf integer
State = {}

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


---@class MarkBlock
---@field

---@param bufnr integer
---@return State
State.new = function(bufnr)
	vim.api.nvim_buf_clear_namespace(bufnr, ns_id, 0, -1)
	local state = setmetatable({ buf = bufnr }, { __index = State })

	local query = vim.treesitter.query.parse('markdown_inline', [[ (latex_block) @latex ]])
	local tree = vim.treesitter.get_parser(bufnr):parse()[1]

	for _, match, _ in query:iter_matches(tree:root(), bufnr, 0, -1, { all = true }) do
		for _, nodes in pairs(match) do
			for _, node in ipairs(nodes) do
				local start_row, start_col, end_row, end_col = node:range()
				-- FIXME: start_row+1 to skip $$ but will probably break something down the line
				state:create_conceal(start_row + 1, end_row)
			end
		end
	end

	return state
end

function State:create_conceal(start_row, end_row)
	vim.system({ 'mathcha', '-render' }, {
		stdin = vim.api.nvim_buf_get_lines(self.buf, start_row, end_row, false)
	}, function(obj)
		-- TODO: check exit code
		vim.schedule(function()
			local lines = enumerate(vim.gsplit(obj.stdout, '\n', { plain = true }))
			---@type integer?
			local last_replaced_row
			for i, line in lines do
				local row = start_row + i
				last_replaced_row = row
				if row ~= end_row then
					vim.api.nvim_buf_set_extmark(self.buf, ns_id, row, 0, {
						-- invalidate = true,
						virt_text = { { line } },
						virt_text_pos = "overlay",
						virt_text_hide = true,
						conceal = "",
						end_col = 9999,
						strict = false,
					})
				else
					-- squash the rest into one big virt_lines
					---@type string[][][]
					local rest = { { { line } } }
					for _, l in lines do
						table.insert(rest, { { l } })
					end

					vim.api.nvim_buf_set_extmark(self.buf, ns_id, row, 0, {
						-- invalidate = true,
						virt_lines = rest,
						virt_text_pos = "overlay",
						virt_lines_above = true,
						virt_text_hide = true,
						conceal = "",

						-- TODO: probably there's a better way
						end_col = 9999,
						strict = false,
					})
				end
			end

			if last_replaced_row == nil then
				for i = start_row, end_row do
					vim.api.nvim_buf_set_extmark(self.buf, ns_id, i, 0, {
						virt_text_hide = true,
						conceal = "",
						end_col = 9999,
						strict = false,
					})
				end
			elseif last_replaced_row < end_row then
				for i = last_replaced_row + 1, end_row do
					vim.api.nvim_buf_set_extmark(self.buf, ns_id, i, 0, {
						virt_text_hide = true,
						conceal = "",
						end_col = 9999,
						strict = false,
					})
				end
			end
		end)
	end)
end

function State:reset_marks()
	vim.api.nvim_buf_clear_namespace(self.buf, ns_id, 0, -1)
end

function M.attach(bufnr)
	-- local ok, _ = pcall(vim.api.nvim_buf_get_var, bufnr, 'mathcha_state')
	-- if not ok then
	vim.api.nvim_buf_set_var(bufnr, 'mathcha_state', State.new(bufnr))
	-- end
end

function M.instance(bufnr)
	local ok, val = pcall(vim.api.nvim_buf_get_var, bufnr or 0, 'mathcha_state')
	if ok then return val else return nil end
end

M.State = State

return M
