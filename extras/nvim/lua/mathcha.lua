local M = {}
---@type table<integer, State>
M._states = {}

---@class State
---@field buf integer
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


---@class MarkBlock
---@field

---@param bufnr integer
---@return State?
function State.new(bufnr)
	vim.api.nvim_buf_clear_namespace(bufnr, ns_id, 0, -1)
	local state = setmetatable({ buf = bufnr }, { __index = State })

	local query = vim.treesitter.query.parse('markdown_inline', [[ (latex_block) @latex ]])
	local md_inline_tree = vim.treesitter.get_parser(bufnr, "markdown")
		:children()["markdown_inline"]

	if md_inline_tree == nil then
		return nil
	end

	md_inline_tree:parse(
		true, -- TODO: parse only visible region maybe
		function(err, trees)
			if err ~= nil or trees == nil then
				vim.notify(("failed to parse buffer %d: %s"):format(bufnr, err))
				return
			end

			for i, tree in pairs(trees) do
				print("i", i, tree:root():sexpr())
				for _, match, _ in query:iter_matches(tree:root(), bufnr, 0, -1) do
					for _, nodes in pairs(match) do
						for _, node in ipairs(nodes) do
							print("1")
							local start_row, start_col, end_row, end_col = node:range()
							-- FIXME: start_row+1 to skip $$ but will probably break something down the line
							state:create_conceal(start_row + 1, end_row)
						end
					end
				end
			end
		end
	)


	return state
end

function State:create_conceal(start_row, end_row)
	assert(start_row < end_row,
		string.format(
			"invariant (start_row < end_row) violated: start_row=%d end_row=%d",
			start_row, end_row
		))

	vim.system({ 'mathcha', '-render' }, {
		stdin = vim.api.nvim_buf_get_lines(self.buf, start_row, end_row, false)
	}, function(obj)
		-- TODO: check exit code
		if obj.code ~= 0 then
			vim.notify("mathcha failed: " .. obj.stderr, vim.log.levels.ERROR)
		end
		vim.schedule(function()
			local lines = enumerate(vim.gsplit(obj.stdout, '\n', { plain = true }))
			vim.print(lines)
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
						-- why -1 no worky??
						end_col = math.huge,
						strict = true,
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
	local buf = bufnr or vim.fn.bufnr()
	if buf == -1 then
		error("invalid bufnr " .. tostring(bufnr))
	end
	-- TODO: cleanup state on buf delete
	if not M._states[buf] then
		M._states[buf] = State.new(buf)
	end
	print("attached to ", buf)
end

function M.instance(bufnr)
	return M._states[vim.fn.bufnr(bufnr or 0)]
end

M.State = State

return M
