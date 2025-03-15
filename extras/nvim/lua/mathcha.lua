local M = {}

---@class State
---@field buf integer
State = {}

local ns_id = vim.api.nvim_create_namespace("mathcha")

---@param bufnr integer
---@return State, string error
State.new = function(bufnr)
	local state = setmetatable({ buf = bufnr }, { __index = State })
	local ts = vim.treesitter

	local query = vim.treesitter.query.parse('markdown_inline', [[ (latex_block) @latex ]])
	local tree = vim.treesitter.get_parser(bufnr):parse()[1]

	vim.print('tree', tree, tree:root())
	for _, match, _ in query:iter_matches(tree:root(), bufnr, 0, -1, { all = true }) do
		print(1)
		for _, nodes in pairs(match) do
			for _, node in ipairs(nodes) do
				local start_row, start_col, end_row, end_col = node:range()
				state:create_conceal(start_row, start_col, end_row, end_col)
			end
		end
	end
	print(4)

	return state
end

function State.create_conceal(self, start_row, start_col, end_row, end_col)
	print(2)
	vim.system({ 'mathcha', '-render' }, {
		stdin = vim.api.nvim_buf_get_lines(self.buf, start_row, end_row, false)
	}, function(obj)
		print(3)
		-- TODO: check exit code
		vim.schedule(function()
			local lines = vim.tbl_map(function(x)
				return { { x } }
			end, vim.split(obj.stdout, '\n'))
			vim.api.nvim_buf_set_extmark(self.buf, ns_id, start_row, start_col, {
				end_row = end_row,
				end_col = end_col,
				-- invalidate = true,
				virt_lines = lines,
				virt_text_pos = "overlay",
			})
		end)
	end)
end

function M.attach(bufnr)
	-- local ok, _ = pcall(vim.api.nvim_buf_get_var, bufnr, 'mathcha_state')
	-- if not ok then
	vim.api.nvim_buf_set_var(bufnr, 'mathcha_state', State.new(bufnr))
	-- end
end

return M
