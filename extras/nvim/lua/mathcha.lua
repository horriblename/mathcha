local M = {}
---@type table<integer, State>
M._states = {}

---@alias NodeID string

---@class NodeBinding
---@field ext_ids integer[]
---@field running vim.SystemObj?
---@field needs_rerun {start: integer, end_: integer}? used to indicate the running job is outdated and should be rerun

---@class State
---@field buf integer
---@field jobs table<NodeID, NodeBinding>
---@field query vim.treesitter.Query
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
	local state = setmetatable({ buf = bufnr, jobs = {} }, { __index = State })

	state.query = vim.treesitter.query.parse('markdown_inline', [[ (latex_block) @latex ]])
	local md_inline_tree = vim.treesitter.get_parser(bufnr, "markdown")
		:children()["markdown_inline"]

	if md_inline_tree == nil then
		return nil
	end

	md_inline_tree:register_cbs({
		on_changedtree = function(_, tree)
			vim.schedule(function()
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
		assert(tree.root and tree.copy)
		for _, match, _ in self.query:iter_matches(tree:root(), self.buf, 0, -1) do
			for _, nodes in pairs(match) do
				for _, node in ipairs(nodes) do
					local start_row, _, end_row, _ = node:range()
					self:update_conceal(node:id(), start_row + 1, end_row)
				end
			end
		end
	end
end

---@param node_id NodeID
---@param start_row integer
---@param end_row integer
function State:update_conceal(node_id, start_row, end_row)
	assert(start_row < end_row,
		string.format(
			"invariant (start_row < end_row) violated: start_row=%d end_row=%d",
			start_row, end_row
		))

	if self.jobs[node_id] and self.jobs[node_id].running then
		self.jobs[node_id].needs_rerun = { start = start_row, end_ = end_row }
		return
	end

	self.jobs[node_id] = { ext_ids = {} }

	self.jobs[node_id].running = vim.system({ './mathcha', '-render' }, {
		stdin = vim.api.nvim_buf_get_lines(self.buf, start_row, end_row, false)
	}, function(obj)
		-- TODO: check exit code
		if obj.code ~= 0 then
			vim.notify("mathcha failed: " .. obj.stderr, vim.log.levels.ERROR)
		end
		vim.schedule(function()
			local lines = enumerate(vim.gsplit(obj.stdout, '\n', { plain = true }))
			local binding = self.jobs[node_id]
			if binding == nil then
				-- deleted maybe? idk
				return
			end

			for _, x in ipairs(binding.ext_ids) do
				assert(vim.api.nvim_buf_del_extmark(self.buf, ns_id, x))
			end

			local extmarks = {}
			---@type integer?
			local last_replaced_row
			for i, line in lines do
				local row = start_row + i
				last_replaced_row = row
				if row ~= end_row then
					table.insert(extmarks, vim.api.nvim_buf_set_extmark(self.buf, ns_id, row, 0, {
						invalidate = true,
						virt_text = { { line } },
						virt_text_pos = "overlay",
						virt_text_hide = true,
						conceal = "",
						-- why -1 no worky??
						end_col = 999999,
						strict = false,
					}))
				else
					-- squash the rest into one big virt_lines
					---@type string[][][]
					local rest = { { { line } } }
					for _, l in lines do
						table.insert(rest, { { l } })
					end

					table.insert(extmarks, vim.api.nvim_buf_set_extmark(self.buf, ns_id, row, 0, {
						invalidate = true,
						virt_lines = rest,
						virt_text_pos = "overlay",
						virt_lines_above = true,
						virt_text_hide = true,
						conceal = "",

						end_col = 999999,
						strict = false,
					}))
				end
			end

			if last_replaced_row == nil then
				for i = start_row, end_row do
					table.insert(extmarks, vim.api.nvim_buf_set_extmark(self.buf, ns_id, i, 0, {
						virt_text_hide = true,
						conceal = "",
						end_col = 999999,
						strict = false,
					}))
				end
			elseif last_replaced_row < end_row then
				for i = last_replaced_row + 1, end_row do
					table.insert(extmarks, vim.api.nvim_buf_set_extmark(self.buf, ns_id, i, 0, {
						virt_text_hide = true,
						conceal = "",
						end_col = 9999,
						strict = false,
					}))
				end
			end

			binding.ext_ids = extmarks
			binding.running = nil
			local params = binding.needs_rerun
			if params then
				binding.needs_rerun = nil
				self:update_conceal(node_id, params.start, params.end_)
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
