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


---@param bufnr integer
---@return State?
function State.new(bufnr)
	vim.api.nvim_buf_clear_namespace(bufnr, ns_id, 0, -1)
	local state = setmetatable({
		buf = bufnr,
		jobs = {},
		pending = {},
	}, { __index = State })

	-- NOTE: we match the whole block including `$$` because otherwise we get each line separately
	-- for some reason
	state.query = vim.treesitter.query.parse('markdown_inline', [[ (latex_block) @latex ]])
	local md_inline_tree = vim.treesitter.get_parser(bufnr, "markdown")
		:children()["markdown_inline"]

	if md_inline_tree == nil then
		return nil
	end

	md_inline_tree:register_cbs({
		on_changedtree = function(_, tree)
			vim.schedule(function()
				-- HACK: idk what the index should be
				state:_parse_and_render({ [1] = tree })
			end)
		end,
		on_detach = function()
			vim.api.nvim_buf_clear_namespace(bufnr, ns_id, 0, -1)
			M._states[bufnr] = nil
		end
	}, true)

	md_inline_tree:parse(false, function(err, trees)
		if err ~= nil or trees == nil then
			vim.notify_once(
				("failed to parse markdown_inline: %s"):format(bufnr, err),
				vim.log.levels.ERROR
			)
			return
		end
		state:_parse_and_render(trees)
	end)
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
	local node_id = node:id()
	start_row = start_row + 1

	self.jobs[node_id] = { ext_ids = {} }

	self.running = vim.system({ './mathcha', '-render' }, {
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

function M.attach(bufnr)
	local buf = bufnr or vim.fn.bufnr()
	if buf == -1 then
		error("invalid bufnr " .. tostring(bufnr))
	end
	-- TODO: cleanup state on buf delete
	if not M._states[buf] then
		M._states[buf] = State.new(buf)
	end
end

function M.instance(bufnr)
	return M._states[vim.fn.bufnr(bufnr or 0)]
end

-- for testing
function M.unload()
	for buf, _ in pairs(M._states) do
		vim.treesitter.stop(buf)
		vim.api.nvim_buf_clear_namespace(buf, ns_id, 0, -1)
	end
end

M.State = State

return M
