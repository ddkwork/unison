// Copyright ©2021-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package unison

import (
	"time"

	"github.com/ddkwork/unison/enums/paintstyle"

	"github.com/ddkwork/golibrary/mylog"
	"github.com/ddkwork/toolbox/xmath"
	"github.com/ddkwork/toolbox/xmath/geom"
	"github.com/google/uuid"
)

var zeroUUID = uuid.UUID{}

// TableDragData holds the data from a table row drag.
type TableDragData[T TableRowConstraint[T]] struct {
	Table *Table[T]
	Rows  []T
}

// ColumnInfo holds column information.
type ColumnInfo struct {
	ID          int
	Current     float32
	Minimum     float32
	Maximum     float32
	AutoMinimum float32
	AutoMaximum float32
}

type tableCache[T TableRowConstraint[T]] struct {
	row    T
	parent int
	depth  int
	height float32
}

type tableHitRect struct {
	Rect
	handler func()
}

// DefaultTableTheme holds the default TableTheme values for Tables. Modifying this data will not alter existing Tables,
// but will alter any Tables created in the future.
var DefaultTableTheme = TableTheme{
	BackgroundInk:          ContentColor,
	OnBackgroundInk:        OnContentColor,
	BandingInk:             BandingColor,
	OnBandingInk:           OnBandingColor,
	InteriorDividerInk:     InteriorDividerColor,
	SelectionInk:           SelectionColor,
	OnSelectionInk:         OnSelectionColor,
	InactiveSelectionInk:   InactiveSelectionColor,
	OnInactiveSelectionInk: OnInactiveSelectionColor,
	IndirectSelectionInk:   IndirectSelectionColor,
	OnIndirectSelectionInk: OnIndirectSelectionColor,
	Padding:                NewUniformInsets(4),
	HierarchyIndent:        16,
	MinimumRowHeight:       16,
	ColumnResizeSlop:       4,
	ShowRowDivider:         true,
	ShowColumnDivider:      true,
}

// TableTheme holds theming data for a Table.
type TableTheme struct {
	BackgroundInk          Ink
	OnBackgroundInk        Ink
	BandingInk             Ink
	OnBandingInk           Ink
	InteriorDividerInk     Ink
	SelectionInk           Ink
	OnSelectionInk         Ink
	InactiveSelectionInk   Ink
	OnInactiveSelectionInk Ink
	IndirectSelectionInk   Ink
	OnIndirectSelectionInk Ink
	Padding                Insets
	HierarchyColumnID      int
	HierarchyIndent        float32
	MinimumRowHeight       float32
	ColumnResizeSlop       float32
	ShowRowDivider         bool
	ShowColumnDivider      bool
}

// Table provides a control that can display data in columns and rows.
type Table[T TableRowConstraint[T]] struct {
	Panel
	TableTheme
	SelectionChangedCallback func()
	DoubleClickCallback      func()
	DragRemovedRowsCallback  func() // Called whenever a drag removes one or more rows from a model, but only if the source and destination tables were different.
	DropOccurredCallback     func() // Called whenever a drop occurs that modifies the model.
	Columns                  []ColumnInfo
	Model                    TableModel[T]
	filteredRows             []T // Note that we use the difference between nil and an empty slice here
	header                   *TableHeader[T]
	selMap                   map[uuid.UUID]bool
	selAnchor                uuid.UUID
	lastSel                  uuid.UUID
	hitRects                 []tableHitRect
	rowCache                 []tableCache[T]
	lastMouseEnterCellPanel  *Panel
	lastMouseDownCellPanel   *Panel
	interactionRow           int
	interactionColumn        int
	lastMouseMotionRow       int
	lastMouseMotionColumn    int
	startRow                 int
	endBeforeRow             int
	columnResizeStart        float32
	columnResizeBase         float32
	columnResizeOverhead     float32
	PreventUserColumnResize  bool
	awaitingSizeColumnsToFit bool
	awaitingSyncToModel      bool
	selNeedsPrune            bool
	wasDragged               bool
	dividerDrag              bool
}

// NewTable creates a new Table control.
func NewTable[T TableRowConstraint[T]](model TableModel[T]) *Table[T] {
	t := &Table[T]{
		TableTheme:            DefaultTableTheme,
		Model:                 model,
		selMap:                make(map[uuid.UUID]bool),
		interactionRow:        -1,
		interactionColumn:     -1,
		lastMouseMotionRow:    -1,
		lastMouseMotionColumn: -1,
	}
	t.Self = t
	t.SetFocusable(true)
	t.SetSizer(t.DefaultSizes)
	t.GainedFocusCallback = t.DefaultFocusGained
	t.DrawCallback = t.DefaultDraw
	t.UpdateCursorCallback = t.DefaultUpdateCursorCallback
	t.UpdateTooltipCallback = t.DefaultUpdateTooltipCallback
	t.MouseMoveCallback = t.DefaultMouseMove
	t.MouseDownCallback = t.DefaultMouseDown
	t.MouseDragCallback = t.DefaultMouseDrag
	t.MouseUpCallback = t.DefaultMouseUp
	t.MouseEnterCallback = t.DefaultMouseEnter
	t.MouseExitCallback = t.DefaultMouseExit
	t.KeyDownCallback = t.DefaultKeyDown
	t.InstallCmdHandlers(SelectAllItemID, AlwaysEnabled, func(_ any) { t.SelectAll() })
	t.wasDragged = false
	return t
}

// ColumnIndexForID returns the column index with the given ID, or -1 if not found.
func (t *Table[T]) ColumnIndexForID(id int) int {
	for i, c := range t.Columns {
		if c.ID == id {
			return i
		}
	}
	return -1
}

// SetDrawRowRange sets a restricted range for sizing and drawing the table. This is intended primarily to be able to
// draw different sections of the table on separate pages of a display and should not be used for anything requiring
// interactivity.
func (t *Table[T]) SetDrawRowRange(start, endBefore int) {
	t.startRow = start
	t.endBeforeRow = endBefore
}

// ClearDrawRowRange clears any restricted range for sizing and drawing the table.
func (t *Table[T]) ClearDrawRowRange() {
	t.startRow = 0
	t.endBeforeRow = 0
}

// CurrentDrawRowRange returns the range of rows that are considered for sizing and drawing.
func (t *Table[T]) CurrentDrawRowRange() (start, endBefore int) {
	if t.startRow < t.endBeforeRow && t.startRow >= 0 && t.endBeforeRow <= len(t.rowCache) {
		return t.startRow, t.endBeforeRow
	}
	return 0, len(t.rowCache)
}

// DefaultDraw provides the default drawing.
func (t *Table[T]) DefaultDraw1(canvas *Canvas, dirty Rect) {
	selectionInk := t.SelectionInk
	if !t.Focused() {
		selectionInk = t.InactiveSelectionInk
	}

	canvas.DrawRect(dirty, t.BackgroundInk.Paint(canvas, dirty, paintstyle.Fill))

	var insets Insets
	if border := t.Border(); border != nil {
		insets = border.Insets()
	}

	var firstCol int
	x := insets.Left
	for i := range t.Columns {
		x1 := x + t.Columns[i].Current
		if t.ShowColumnDivider {
			x1++
		}
		if x1 >= dirty.X {
			break
		}
		x = x1
		firstCol = i + 1
	}

	startRow, endBeforeRow := t.CurrentDrawRowRange()
	y := insets.Top
	for i := startRow; i < endBeforeRow; i++ {
		y1 := y + t.rowCache[i].height
		if t.ShowRowDivider {
			y1++
		}
		if y1 >= dirty.Y {
			break
		}
		y = y1
		startRow = i + 1
	}

	lastY := dirty.Bottom()
	rect := dirty
	rect.Y = y
	for r := startRow; r < endBeforeRow && rect.Y < lastY; r++ {
		rect.Height = t.rowCache[r].height
		if t.IsRowOrAnyParentSelected(r) {
			if t.IsRowSelected(r) {
				canvas.DrawRect(rect, selectionInk.Paint(canvas, rect, paintstyle.Fill))
			} else {
				canvas.DrawRect(rect, t.IndirectSelectionInk.Paint(canvas, rect, paintstyle.Fill))
			}
		} else if r%2 == 1 {
			canvas.DrawRect(rect, t.BandingInk.Paint(canvas, rect, paintstyle.Fill))
		}
		rect.Y += t.rowCache[r].height
		if t.ShowRowDivider && r != endBeforeRow-1 {
			rect.Height = 1
			canvas.DrawRect(rect, t.InteriorDividerInk.Paint(canvas, rect, paintstyle.Fill))
			rect.Y++
		}
	}

	if t.ShowColumnDivider {
		rect = dirty
		rect.X = x
		rect.Width = 1
		for c := firstCol; c < len(t.Columns)-1; c++ {
			rect.X += t.Columns[c].Current
			canvas.DrawRect(rect, t.InteriorDividerInk.Paint(canvas, rect, paintstyle.Fill))
			rect.X++
		}
	}

	rect = dirty
	rect.Y = y
	lastX := dirty.Right()
	t.hitRects = nil
	for r := startRow; r < endBeforeRow && rect.Y < lastY; r++ {
		rect.X = x
		rect.Height = t.rowCache[r].height
		for c := firstCol; c < len(t.Columns) && rect.X < lastX; c++ {
			fg, bg, selected, indirectlySelected, focused := t.cellParams(r, c)
			rect.Width = t.Columns[c].Current
			cellRect := rect
			cellRect.Inset(t.Padding)
			row := t.rowCache[r].row
			if t.Columns[c].ID == t.HierarchyColumnID {
				if row.CanHaveChildren() {
					const disclosureIndent = 2
					disclosureSize := min(t.HierarchyIndent, t.MinimumRowHeight) - disclosureIndent*2
					canvas.Save()
					left := cellRect.X + t.HierarchyIndent*float32(t.rowCache[r].depth) + disclosureIndent
					top := cellRect.Y + (t.MinimumRowHeight-disclosureSize)/2
					t.hitRects = append(t.hitRects, t.newTableHitRect(NewRect(left, top, disclosureSize,
						disclosureSize), row))
					canvas.Translate(left, top)
					if row.IsOpen() {
						offset := disclosureSize / 2
						canvas.Translate(offset, offset)
						canvas.Rotate(90)
						canvas.Translate(-offset, -offset)
					}
					canvas.DrawPath(CircledChevronRightSVG.PathForSize(NewSize(disclosureSize, disclosureSize)),
						fg.Paint(canvas, cellRect, paintstyle.Fill))
					canvas.Restore()
				}
				indent := t.HierarchyIndent*float32(t.rowCache[r].depth+1) + t.Padding.Left
				cellRect.X += indent
				cellRect.Width -= indent
			}
			cell := row.ColumnCell(r, c, fg, bg, selected, indirectlySelected, focused).AsPanel()
			t.installCell(cell, cellRect)
			canvas.Save()
			canvas.Translate(cellRect.X, cellRect.Y)
			cellRect.X = 0
			cellRect.Y = 0
			cell.Draw(canvas, cellRect)
			t.uninstallCell(cell)
			canvas.Restore()
			rect.X += t.Columns[c].Current
			if t.ShowColumnDivider {
				rect.X++
			}
		}
		rect.Y += t.rowCache[r].height
		if t.ShowRowDivider {
			rect.Y++
		}
	}
}

// DefaultDraw 提供默认绘制方法
func (t *Table[T]) DefaultDraw(canvas *Canvas, dirty Rect) {
	selectionInk := t.SelectionInk // 获取选中状态的墨水颜色
	if !t.Focused() {              // 如果表格没有获得焦点
		selectionInk = t.InactiveSelectionInk // 使用未激活的选中墨水颜色
	}

	canvas.DrawRect(dirty, t.BackgroundInk.Paint(canvas, dirty, paintstyle.Fill)) // 绘制背景矩形

	var insets Insets
	if border := t.Border(); border != nil { // 获取边框的内边距
		insets = border.Insets()
	}

	var firstCol int
	x := insets.Left           // 设置起始位置为左内边距
	for i := range t.Columns { // 遍历所有列
		x1 := x + t.Columns[i].Current // 计算当前列的右侧位置
		if t.ShowColumnDivider {       // 如果显示列分隔符
			x1++ // 列分隔符占用一个位置
		}
		if x1 >= dirty.X { // 如果当前列的右侧位置已经超出可绘制区域
			break // 结束循环
		}
		x = x1           // 更新 x 位置为当前列的右侧位置
		firstCol = i + 1 // 设置首列索引
	}

	startRow, endBeforeRow := t.CurrentDrawRowRange() // 获取当前可绘制的行范围
	y := insets.Top                                   // 设置起始 y 位置为上内边距
	for i := startRow; i < endBeforeRow; i++ {        // 遍历可绘制的行
		y1 := y + t.rowCache[i].height // 计算当前行的底部 y 位置
		if t.ShowRowDivider {          // 如果显示行分隔符
			y1++ // 行分隔符占用一个位置
		}
		if y1 >= dirty.Y { // 如果当前行的底部位置已经超出可绘制区域
			break // 结束循环
		}
		y = y1           // 更新 y 位置为当前行的底部位置
		startRow = i + 1 // 更新起始行索引
	}

	lastY := dirty.Bottom()                                      // 获取可绘制区域的底部 y 坐标
	rect := dirty                                                // 创建矩形用于绘制
	rect.Y = y                                                   // 设置矩形的顶部 y 坐标
	for r := startRow; r < endBeforeRow && rect.Y < lastY; r++ { // 遍历可绘制的行
		rect.Height = t.rowCache[r].height // 设置矩形的高度为当前行的高度
		if t.IsRowOrAnyParentSelected(r) { // 如果当前行或其父级被选中
			if t.IsRowSelected(r) { // 如果当前行被选中
				canvas.DrawRect(rect, selectionInk.Paint(canvas, rect, paintstyle.Fill)) // 绘制选中行的背景
			} else {
				canvas.DrawRect(rect, t.IndirectSelectionInk.Paint(canvas, rect, paintstyle.Fill)) // 绘制间接选中行的背景
			}
		} else if r%2 == 1 { // 如果当前行是奇数行
			canvas.DrawRect(rect, t.BandingInk.Paint(canvas, rect, paintstyle.Fill)) // 绘制带状背景
		}
		rect.Y += t.rowCache[r].height               // 更新矩形的顶部 y 坐标为下一行的高度
		if t.ShowRowDivider && r != endBeforeRow-1 { // 如果显示行分隔符并且不是最后一行
			rect.Height = 1                                                                  // 设置行分隔符的高度为 1
			canvas.DrawRect(rect, t.InteriorDividerInk.Paint(canvas, rect, paintstyle.Fill)) // 绘制行分隔符
			rect.Y++                                                                         // 更新矩形的顶部 y 坐标
		}
	}

	if t.ShowColumnDivider { // 如果显示列分隔符
		rect = dirty                                   // 重置矩形为可绘制区域
		rect.X = x                                     // 设置矩形的左侧 x 坐标
		rect.Width = 1                                 // 设置矩形的宽度为 1
		for c := firstCol; c < len(t.Columns)-1; c++ { // 遍历列
			rect.X += t.Columns[c].Current                                                   // 更新矩形的左侧 x 坐标
			canvas.DrawRect(rect, t.InteriorDividerInk.Paint(canvas, rect, paintstyle.Fill)) // 绘制列分隔符
			rect.X++                                                                         // 更新矩形的左侧 x 坐标
		}
	}

	rect = dirty                                                 // 重置矩形为可绘制区域
	rect.Y = y                                                   // 更新矩形的顶部 y 坐标
	lastX := dirty.Right()                                       // 获取可绘制区域的右侧 x 坐标
	t.hitRects = nil                                             // 清空 hitRects
	for r := startRow; r < endBeforeRow && rect.Y < lastY; r++ { // 遍历可绘制的行
		rect.X = x                                                     // 设置矩形的左侧 x 坐标
		rect.Height = t.rowCache[r].height                             // 设置矩形的高度为当前行的高度
		for c := firstCol; c < len(t.Columns) && rect.X < lastX; c++ { // 遍历可绘制的列
			fg, bg, selected, indirectlySelected, focused := t.cellParams(r, c) // 获取当前单元格的参数
			rect.Width = t.Columns[c].Current                                   // 设置矩形的宽度为当前列的宽度
			cellRect := rect                                                    // 保存当前矩形
			cellRect.Inset(t.Padding)                                           // 设置单元格的内边距
			row := t.rowCache[r].row                                            // 获取当前行的数据
			if t.Columns[c].ID == t.HierarchyColumnID {                         // 如果当前列是层级列
				if row.CanHaveChildren() { // 如果当前行可以有子项
					const disclosureIndent = 2                                                                                  // 设置展开图标的边距
					disclosureSize := min(t.HierarchyIndent, t.MinimumRowHeight) - disclosureIndent*2                           // 设置展开图标的大小
					canvas.Save()                                                                                               // 保存当前画布状态
					left := cellRect.X + t.HierarchyIndent*float32(t.rowCache[r].depth) + disclosureIndent                      // 计算展开图标的左侧 x 坐标
					top := cellRect.Y + (t.MinimumRowHeight-disclosureSize)/2                                                   // 计算展开图标的顶部 y 坐标
					t.hitRects = append(t.hitRects, t.newTableHitRect(NewRect(left, top, disclosureSize, disclosureSize), row)) // 添加到 hitRects
					canvas.Translate(left, top)                                                                                 // 移动画布到展开图标的位置
					if row.IsOpen() {                                                                                           // 如果当前行是展开状态
						offset := disclosureSize / 2       // 计算偏移量
						canvas.Translate(offset, offset)   // 移动画布
						canvas.Rotate(90)                  // 旋转画布 90 度
						canvas.Translate(-offset, -offset) // 移动回原来位置
					}
					canvas.DrawPath(CircledChevronRightSVG.PathForSize(NewSize(disclosureSize, disclosureSize)), fg.Paint(canvas, cellRect, paintstyle.Fill)) // 绘制展开图标
					canvas.Restore()                                                                                                                          // 恢复画布状态
				}
				indent := t.HierarchyIndent*float32(t.rowCache[r].depth+1) + t.Padding.Left // 计算缩进
				cellRect.X += indent                                                        // 更新单元格的左侧 x 坐标
				cellRect.Width -= indent                                                    // 更新单元格的宽度
			}
			cell := row.ColumnCell(r, c, fg, bg, selected, indirectlySelected, focused).AsPanel() // 获取当前单元格的面板
			t.installCell(cell, cellRect)                                                         // 安装单元格
			canvas.Save()                                                                         // 保存当前画布状态
			canvas.Translate(cellRect.X, cellRect.Y)                                              // 移动画布到单元格的位置
			cellRect.X = 0                                                                        // 重置单元格矩形的位置
			cellRect.Y = 0                                                                        // 重置单元格矩形的位置
			cell.Draw(canvas, cellRect)                                                           // 绘制单元格
			t.uninstallCell(cell)                                                                 // 卸载单元格
			canvas.Restore()                                                                      // 恢复画布状态
			rect.X += t.Columns[c].Current                                                        // 更新矩形的左侧 x 坐标
			if t.ShowColumnDivider {                                                              // 如果显示列分隔符
				rect.X++ // 更新矩形的左侧 x 坐标
			}
		}
		rect.Y += t.rowCache[r].height // 更新矩形的顶部 y 坐标
		if t.ShowRowDivider {          // 如果显示行分隔符
			rect.Y++ // 更新矩形的顶部 y 坐标
		}
	}
}

func (t *Table[T]) cellParams(row, _ int) (fg, bg Ink, selected, indirectlySelected, focused bool) {
	focused = t.Focused()
	selected = t.IsRowSelected(row)
	indirectlySelected = !selected && t.IsRowOrAnyParentSelected(row)
	switch {
	case selected && focused:
		fg = t.OnSelectionInk
		bg = t.SelectionInk
	case selected:
		fg = t.OnInactiveSelectionInk
		bg = t.InactiveSelectionInk
	case indirectlySelected:
		fg = t.OnIndirectSelectionInk
		bg = t.IndirectSelectionInk
	case row%2 == 1:
		fg = t.OnBandingInk
		bg = t.BandingInk
	default:
		fg = t.OnBackgroundInk
		bg = t.BackgroundInk
	}
	return fg, bg, selected, indirectlySelected, focused
}

func (t *Table[T]) cell(row, col int) *Panel {
	fg, bg, selected, indirectlySelected, focused := t.cellParams(row, col)
	return t.rowCache[row].row.ColumnCell(row, col, fg, bg, selected, indirectlySelected, focused).AsPanel()
}

func (t *Table[T]) installCell(cell *Panel, frame Rect) {
	cell.SetFrameRect(frame)
	cell.ValidateLayout()
	cell.parent = t.AsPanel()
}

func (t *Table[T]) uninstallCell(cell *Panel) {
	cell.parent = nil
}

// RowHeights returns the heights of each row.
func (t *Table[T]) RowHeights() []float32 {
	heights := make([]float32, len(t.rowCache))
	for i := range t.rowCache {
		heights[i] = t.rowCache[i].height
	}
	return heights
}

// OverRow returns the row index that the y coordinate is over, or -1 if it isn't over any row.
func (t *Table[T]) OverRow(y float32) int {
	var insets Insets
	if border := t.Border(); border != nil {
		insets = border.Insets()
	}
	end := insets.Top
	for i := range t.rowCache {
		start := end
		end += t.rowCache[i].height
		if t.ShowRowDivider {
			end++
		}
		if y >= start && y < end {
			return i
		}
	}
	return -1
}

// OverColumn returns the column index that the x coordinate is over, or -1 if it isn't over any column.
func (t *Table[T]) OverColumn(x float32) int {
	var insets Insets
	if border := t.Border(); border != nil {
		insets = border.Insets()
	}
	end := insets.Left
	for i := range t.Columns {
		start := end
		end += t.Columns[i].Current
		if t.ShowColumnDivider {
			end++
		}
		if x >= start && x < end {
			return i
		}
	}
	return -1
}

// OverColumnDivider returns the column index of the column divider that the x coordinate is over, or -1 if it isn't
// over any column divider.
func (t *Table[T]) OverColumnDivider(x float32) int {
	if len(t.Columns) < 2 {
		return -1
	}
	var insets Insets
	if border := t.Border(); border != nil {
		insets = border.Insets()
	}
	pos := insets.Left
	for i := range t.Columns[:len(t.Columns)-1] {
		pos += t.Columns[i].Current
		if t.ShowColumnDivider {
			pos++
		}
		if xmath.Abs(pos-x) < t.ColumnResizeSlop {
			return i
		}
	}
	return -1
}

// CellWidth returns the current width of a given cell.
func (t *Table[T]) CellWidth1(row, col int) float32 {
	if row < 0 || col < 0 || row >= len(t.rowCache) || col >= len(t.Columns) {
		return 0
	}
	width := t.Columns[col].Current - (t.Padding.Left + t.Padding.Right)
	if t.Columns[col].ID == t.HierarchyColumnID {
		width -= t.HierarchyIndent*float32(t.rowCache[row].depth+1) + t.Padding.Left
	}
	return width
}

// CellWidth 返回给定单元格的当前宽度
func (t *Table[T]) CellWidth(row, col int) float32 {
	// 检查行和列索引是否有效
	if row < 0 || col < 0 || row >= len(t.rowCache) || col >= len(t.Columns) {
		return 0 // 如果无效，返回 0
	}

	// 获取当前列的宽度，减去左和右的内边距
	width := t.Columns[col].Current - (t.Padding.Left + t.Padding.Right)

	// 如果当前列是层级列
	if t.Columns[col].ID == t.HierarchyColumnID {
		// 根据行的深度计算额外的缩进量
		width -= t.HierarchyIndent*float32(t.rowCache[row].depth+1) + t.Padding.Left
	}

	// 返回计算后的单元格宽度
	return width
}

// ColumnEdges 返回给定列的左侧和右侧的 x 坐标
func (t *Table[T]) ColumnEdges(col int) (left, right float32) {
	// 检查列索引是否有效
	if col < 0 || col >= len(t.Columns) {
		return 0, 0 // 如果无效，返回 0
	}

	var insets Insets
	if border := t.Border(); border != nil { // 获取边框的内边距
		insets = border.Insets()
	}

	left = insets.Left         // 设置左侧 x 坐标为左内边距
	for c := 0; c < col; c++ { // 遍历当前列之前的所有列
		left += t.Columns[c].Current // 累加每列的宽度
		if t.ShowColumnDivider {     // 如果显示列分隔符
			left++ // 列分隔符占用一个位置
		}
	}

	right = left + t.Columns[col].Current // 右侧 x 坐标为左侧加上当前列的宽度
	left += t.Padding.Left                // 添加左内边距
	right -= t.Padding.Right              // 减去右内边距

	// 如果当前列是层级列，增加额外的缩进
	if t.Columns[col].ID == t.HierarchyColumnID {
		left += t.HierarchyIndent + t.Padding.Left
	}

	// 确保右侧坐标不小于左侧
	if right < left {
		right = left
	}

	return left, right // 返回列的左右边界
}

// CellFrame 返回给定单元格的矩形框架
func (t *Table[T]) CellFrame(row, col int) Rect {
	// 检查行和列索引是否有效
	if row < 0 || col < 0 || row >= len(t.rowCache) || col >= len(t.Columns) {
		return Rect{} // 如果无效，返回空矩形
	}

	var insets Insets
	if border := t.Border(); border != nil { // 获取边框的内边距
		insets = border.Insets()
	}

	x := insets.Left           // 设置左侧 x 坐标为左内边距
	for c := 0; c < col; c++ { // 遍历当前列之前的所有列
		x += t.Columns[c].Current // 累加每列的宽度
		if t.ShowColumnDivider {  // 如果显示列分隔符
			x++ // 列分隔符占用一个位置
		}
	}

	y := insets.Top            // 设置顶部 y 坐标为上内边距
	for r := 0; r < row; r++ { // 遍历当前行之前的所有行
		y += t.rowCache[r].height // 累加每行的高度
		if t.ShowRowDivider {     // 如果显示行分隔符
			y++ // 行分隔符占用一个位置
		}
	}

	// 创建并返回单元格的矩形
	rect := NewRect(x, y, t.Columns[col].Current, t.rowCache[row].height)
	rect.Inset(t.Padding) // 设置单元格的内边距
	// 如果当前列是层级列，添加缩进
	if t.Columns[col].ID == t.HierarchyColumnID {
		indent := t.HierarchyIndent*float32(t.rowCache[row].depth+1) + t.Padding.Left // 层级缩进*深度+1+左内边距
		rect.X += indent                                                              // 更新矩形的左侧 x 坐标
		rect.Width -= indent                                                          // 更新矩形的宽度
		if rect.Width < 1 {                                                           // 确保宽度不小于 1
			rect.Width = 1
		}
	}
	return rect // 返回单元格的矩形框架
}

// RowFrame 返回给定行的矩形框架
func (t *Table[T]) RowFrame(row int) Rect {
	// 检查行索引是否有效
	if row < 0 || row >= len(t.rowCache) {
		return Rect{} // 如果无效，返回空矩形
	}

	rect := t.ContentRect(false) // 获取内容区域的矩形
	for i := 0; i < row; i++ {   // 遍历当前行之前的所有行
		rect.Y += t.rowCache[i].height // 累加每行的高度
		if t.ShowRowDivider {          // 如果显示行分隔符
			rect.Y++ // 行分隔符占用一个位置
		}
	}

	rect.Height = t.rowCache[row].height // 设置矩形的高度为当前行的高度
	return rect                          // 返回行的矩形框架
}

// ColumnEdges returns the x-coordinates of the left and right sides of the column.
func (t *Table[T]) ColumnEdges1(col int) (left, right float32) {
	if col < 0 || col >= len(t.Columns) {
		return 0, 0
	}
	var insets Insets
	if border := t.Border(); border != nil {
		insets = border.Insets()
	}
	left = insets.Left
	for c := 0; c < col; c++ {
		left += t.Columns[c].Current
		if t.ShowColumnDivider {
			left++
		}
	}
	right = left + t.Columns[col].Current
	left += t.Padding.Left
	right -= t.Padding.Right
	if t.Columns[col].ID == t.HierarchyColumnID {
		left += t.HierarchyIndent + t.Padding.Left
	}
	if right < left {
		right = left
	}
	return left, right
}

// CellFrame returns the frame of the given cell.
func (t *Table[T]) CellFrame1(row, col int) Rect {
	if row < 0 || col < 0 || row >= len(t.rowCache) || col >= len(t.Columns) {
		return Rect{}
	}
	var insets Insets
	if border := t.Border(); border != nil {
		insets = border.Insets()
	}
	x := insets.Left
	for c := 0; c < col; c++ {
		x += t.Columns[c].Current
		if t.ShowColumnDivider {
			x++
		}
	}
	y := insets.Top
	for r := 0; r < row; r++ {
		y += t.rowCache[r].height
		if t.ShowRowDivider {
			y++
		}
	}
	rect := NewRect(x, y, t.Columns[col].Current, t.rowCache[row].height)
	rect.Inset(t.Padding)
	if t.Columns[col].ID == t.HierarchyColumnID {
		indent := t.HierarchyIndent*float32(t.rowCache[row].depth+1) + t.Padding.Left
		rect.X += indent
		rect.Width -= indent
		if rect.Width < 1 {
			rect.Width = 1
		}
	}
	return rect
}

// RowFrame returns the frame of the row.
func (t *Table[T]) RowFrame1(row int) Rect {
	if row < 0 || row >= len(t.rowCache) {
		return Rect{}
	}
	rect := t.ContentRect(false)
	for i := 0; i < row; i++ {
		rect.Y += t.rowCache[i].height
		if t.ShowRowDivider {
			rect.Y++
		}
	}
	rect.Height = t.rowCache[row].height
	return rect
}

func (t *Table[T]) newTableHitRect(rect Rect, row T) tableHitRect {
	return tableHitRect{
		Rect: rect,
		handler: func() {
			open := !row.IsOpen()
			row.SetOpen(open)
			t.SyncToModel()
			if !open {
				t.PruneSelectionOfUndisclosedNodes()
			}
		},
	}
}

// DefaultFocusGained provides the default focus gained handling.
func (t *Table[T]) DefaultFocusGained() {
	switch {
	case t.interactionRow != -1:
		t.ScrollRowIntoView(t.interactionRow)
	case t.lastMouseMotionRow != -1:
		t.ScrollRowIntoView(t.lastMouseMotionRow)
	default:
		t.ScrollIntoView()
	}
	t.MarkForRedraw()
}

// DefaultUpdateCursorCallback provides the default cursor update handling.
func (t *Table[T]) DefaultUpdateCursorCallback(where Point) *Cursor {
	if !t.PreventUserColumnResize {
		if over := t.OverColumnDivider(where.X); over != -1 {
			if t.Columns[over].Minimum <= 0 || t.Columns[over].Minimum < t.Columns[over].Maximum {
				return ResizeHorizontalCursor()
			}
		}
	}
	if row := t.OverRow(where.Y); row != -1 {
		if col := t.OverColumn(where.X); col != -1 {
			cell := t.cell(row, col)
			if cell.HasInSelfOrDescendants(func(p *Panel) bool { return p.UpdateCursorCallback != nil }) {
				var cursor *Cursor
				rect := t.CellFrame(row, col)
				t.installCell(cell, rect)
				where.Subtract(rect.Point)
				target := cell.PanelAt(where)
				for target != t.AsPanel() {
					if target.UpdateCursorCallback == nil {
						target = target.parent
					} else {
						mylog.Call(func() { cursor = target.UpdateCursorCallback(cell.PointTo(where, target)) })
						break
					}
				}
				t.uninstallCell(cell)
				return cursor
			}
		}
	}
	return nil
}

// DefaultUpdateTooltipCallback provides the default tooltip update handling.
func (t *Table[T]) DefaultUpdateTooltipCallback(where Point, avoid Rect) Rect {
	if row := t.OverRow(where.Y); row != -1 {
		if col := t.OverColumn(where.X); col != -1 {
			cell := t.cell(row, col)
			if cell.HasInSelfOrDescendants(func(p *Panel) bool { return p.UpdateTooltipCallback != nil || p.Tooltip != nil }) {
				rect := t.CellFrame(row, col)
				t.installCell(cell, rect)
				where.Subtract(rect.Point)
				target := cell.PanelAt(where)
				t.Tooltip = nil
				t.TooltipImmediate = false
				for target != t.AsPanel() {
					avoid = target.RectToRoot(target.ContentRect(true))
					avoid.Align()
					if target.UpdateTooltipCallback != nil {
						mylog.Call(func() { avoid = target.UpdateTooltipCallback(cell.PointTo(where, target), avoid) })
					}
					if target.Tooltip != nil {
						t.Tooltip = target.Tooltip
						t.TooltipImmediate = target.TooltipImmediate
						break
					}
					target = target.parent
				}
				t.uninstallCell(cell)
				return avoid
			}
			if cell.Tooltip != nil {
				t.Tooltip = cell.Tooltip
				t.TooltipImmediate = cell.TooltipImmediate
				avoid = t.RectToRoot(t.CellFrame(row, col))
				avoid.Align()
				return avoid
			}
		}
	}
	t.Tooltip = nil
	return Rect{}
}

// DefaultMouseEnter provides the default mouse enter handling.
func (t *Table[T]) DefaultMouseEnter(where Point, mod Modifiers) bool {
	row := t.OverRow(where.Y)
	col := t.OverColumn(where.X)
	if t.lastMouseMotionRow != row || t.lastMouseMotionColumn != col {
		t.DefaultMouseExit()
		t.lastMouseMotionRow = row
		t.lastMouseMotionColumn = col
	}
	if row != -1 && col != -1 {
		cell := t.cell(row, col)
		rect := t.CellFrame(row, col)
		t.installCell(cell, rect)
		where.Subtract(rect.Point)
		target := cell.PanelAt(where)
		if target != t.lastMouseEnterCellPanel && t.lastMouseEnterCellPanel != nil {
			t.DefaultMouseExit()
			t.lastMouseMotionRow = row
			t.lastMouseMotionColumn = col
		}
		if target.MouseEnterCallback != nil {
			mylog.Call(func() { target.MouseEnterCallback(cell.PointTo(where, target), mod) })
		}
		t.uninstallCell(cell)
		t.lastMouseEnterCellPanel = target
	}
	return true
}

// DefaultMouseMove provides the default mouse move handling.
func (t *Table[T]) DefaultMouseMove(where Point, mod Modifiers) bool {
	t.DefaultMouseEnter(where, mod)
	if t.lastMouseEnterCellPanel != nil {
		row := t.OverRow(where.Y)
		col := t.OverColumn(where.X)
		cell := t.cell(row, col)
		rect := t.CellFrame(row, col)
		t.installCell(cell, rect)
		where.Subtract(rect.Point)
		if target := cell.PanelAt(where); target.MouseMoveCallback != nil {
			mylog.Call(func() { target.MouseMoveCallback(cell.PointTo(where, target), mod) })
		}
		t.uninstallCell(cell)
	}
	return true
}

// DefaultMouseExit provides the default mouse exit handling.
func (t *Table[T]) DefaultMouseExit() bool {
	if t.lastMouseEnterCellPanel != nil && t.lastMouseEnterCellPanel.MouseExitCallback != nil &&
		t.lastMouseMotionColumn != -1 && t.lastMouseMotionRow >= 0 && t.lastMouseMotionRow < len(t.rowCache) {
		cell := t.cell(t.lastMouseMotionRow, t.lastMouseMotionColumn)
		rect := t.CellFrame(t.lastMouseMotionRow, t.lastMouseMotionColumn)
		t.installCell(cell, rect)
		mylog.Call(func() { t.lastMouseEnterCellPanel.MouseExitCallback() })
		t.uninstallCell(cell)
	}
	t.lastMouseEnterCellPanel = nil
	t.lastMouseMotionRow = -1
	t.lastMouseMotionColumn = -1
	return true
}

// DefaultMouseDown provides the default mouse down handling.
func (t *Table[T]) DefaultMouseDown(where Point, button, clickCount int, mod Modifiers) bool {
	if t.Window().InDrag() {
		return false
	}
	t.RequestFocus()
	t.wasDragged = false
	t.dividerDrag = false
	t.lastSel = zeroUUID

	t.interactionRow = -1
	t.interactionColumn = -1
	if button == ButtonLeft {
		if !t.PreventUserColumnResize {
			if over := t.OverColumnDivider(where.X); over != -1 {
				if t.Columns[over].Minimum <= 0 || t.Columns[over].Minimum < t.Columns[over].Maximum {
					if clickCount == 2 {
						t.SizeColumnToFit(over, true)
						t.MarkForRedraw()
						t.Window().UpdateCursorNow()
						return true
					}
					t.interactionColumn = over
					t.columnResizeStart = where.X
					t.columnResizeBase = t.Columns[over].Current
					t.columnResizeOverhead = t.Padding.Left + t.Padding.Right
					if t.Columns[over].ID == t.HierarchyColumnID {
						depth := 0
						for _, cache := range t.rowCache {
							if depth < cache.depth {
								depth = cache.depth
							}
						}
						t.columnResizeOverhead += t.Padding.Left + t.HierarchyIndent*float32(depth+1)
					}
					return true
				}
			}
		}
		for _, one := range t.hitRects {
			if one.ContainsPoint(where) {
				return true
			}
		}
	}
	if row := t.OverRow(where.Y); row != -1 {
		if col := t.OverColumn(where.X); col != -1 {
			cell := t.cell(row, col)
			if cell.HasInSelfOrDescendants(func(p *Panel) bool { return p.MouseDownCallback != nil }) {
				t.interactionRow = row
				t.interactionColumn = col
				rect := t.CellFrame(row, col)
				t.installCell(cell, rect)
				where.Subtract(rect.Point)
				stop := false
				if target := cell.PanelAt(where); target.MouseDownCallback != nil {
					t.lastMouseDownCellPanel = target
					mylog.Call(func() {
						stop = target.MouseDownCallback(cell.PointTo(where, target), button,
							clickCount, mod)
					})
				}
				t.uninstallCell(cell)
				if stop {
					return stop
				}
			}
		}
		rowData := t.rowCache[row].row
		id := rowData.UUID()
		switch {
		case mod&ShiftModifier != 0: // Extend selection from anchor
			selAnchorIndex := -1
			if t.selAnchor != zeroUUID {
				for i, c := range t.rowCache {
					if c.row.UUID() == t.selAnchor {
						selAnchorIndex = i
						break
					}
				}
			}
			if selAnchorIndex != -1 {
				last := max(selAnchorIndex, row)
				for i := min(selAnchorIndex, row); i <= last; i++ {
					t.selMap[t.rowCache[i].row.UUID()] = true
				}
				t.notifyOfSelectionChange()
			} else if !t.selMap[id] { // No anchor, so behave like a regular click
				t.selMap = make(map[uuid.UUID]bool)
				t.selMap[id] = true
				t.selAnchor = id
				t.notifyOfSelectionChange()
			}
		case mod.DiscontiguousSelectionDown(): // Toggle single row
			if t.selMap[id] {
				delete(t.selMap, id)
			} else {
				t.selMap[id] = true
			}
			t.notifyOfSelectionChange()
		case t.selMap[id]: // Sets lastClick so that on mouse up, we can treat a click and click and hold differently
			t.lastSel = id
		default: // If not already selected, replace selection with current row and make it the anchor
			t.selMap = make(map[uuid.UUID]bool)
			t.selMap[id] = true
			t.selAnchor = id
			t.notifyOfSelectionChange()
		}
		t.MarkForRedraw()
		if button == ButtonLeft && clickCount == 2 && t.DoubleClickCallback != nil && len(t.selMap) != 0 {
			mylog.Call(t.DoubleClickCallback)
		}
	}
	return true
}

func (t *Table[T]) notifyOfSelectionChange() {
	if t.SelectionChangedCallback != nil {
		mylog.Call(t.SelectionChangedCallback)
	}
}

// DefaultMouseDrag provides the default mouse drag handling.
func (t *Table[T]) DefaultMouseDrag(where Point, button int, mod Modifiers) bool {
	t.wasDragged = true
	stop := false
	if t.interactionColumn != -1 {
		if t.interactionRow == -1 {
			if button == ButtonLeft && !t.PreventUserColumnResize {
				width := t.columnResizeBase + where.X - t.columnResizeStart
				if width < t.columnResizeOverhead {
					width = t.columnResizeOverhead
				}
				minimum := t.Columns[t.interactionColumn].Minimum
				if minimum > 0 && width < minimum+t.columnResizeOverhead {
					width = minimum + t.columnResizeOverhead
				} else {
					maximum := t.Columns[t.interactionColumn].Maximum
					if maximum > 0 && width > maximum+t.columnResizeOverhead {
						width = maximum + t.columnResizeOverhead
					}
				}
				if t.Columns[t.interactionColumn].Current != width {
					t.Columns[t.interactionColumn].Current = width
					t.EventuallySyncToModel()
					t.MarkForRedraw()
					t.dividerDrag = true
				}
				stop = true
			}
		} else if t.lastMouseDownCellPanel != nil && t.lastMouseDownCellPanel.MouseDragCallback != nil {
			cell := t.cell(t.interactionRow, t.interactionColumn)
			rect := t.CellFrame(t.interactionRow, t.interactionColumn)
			t.installCell(cell, rect)
			where.Subtract(rect.Point)
			mylog.Call(func() {
				stop = t.lastMouseDownCellPanel.MouseDragCallback(cell.PointTo(where, t.lastMouseDownCellPanel), button, mod)
			})
			t.uninstallCell(cell)
		}
	}
	return stop
}

// DefaultMouseUp provides the default mouse up handling.
func (t *Table[T]) DefaultMouseUp(where Point, button int, mod Modifiers) bool {
	stop := false
	if !t.dividerDrag && button == ButtonLeft {
		for _, one := range t.hitRects {
			if one.ContainsPoint(where) {
				one.handler()
				stop = true
				break
			}
		}
	}

	if !t.wasDragged && t.lastSel != zeroUUID {
		t.ClearSelection()
		t.selMap[t.lastSel] = true
		t.selAnchor = t.lastSel
		t.MarkForRedraw()
		t.notifyOfSelectionChange()
	}

	if !stop && t.interactionRow != -1 && t.interactionColumn != -1 && t.lastMouseDownCellPanel != nil &&
		t.lastMouseDownCellPanel.MouseUpCallback != nil {
		cell := t.cell(t.interactionRow, t.interactionColumn)
		rect := t.CellFrame(t.interactionRow, t.interactionColumn)
		t.installCell(cell, rect)
		where.Subtract(rect.Point)
		mylog.Call(func() {
			stop = t.lastMouseDownCellPanel.MouseUpCallback(cell.PointTo(where, t.lastMouseDownCellPanel), button, mod)
		})
		t.uninstallCell(cell)
	}
	t.lastMouseDownCellPanel = nil
	return stop
}

// DefaultKeyDown provides the default key down handling.
func (t *Table[T]) DefaultKeyDown(keyCode KeyCode, mod Modifiers, _ bool) bool {
	if IsControlAction(keyCode, mod) {
		if t.DoubleClickCallback != nil && len(t.selMap) != 0 {
			mylog.Call(t.DoubleClickCallback)
		}
		return true
	}
	switch keyCode {
	case KeyLeft:
		if t.HasSelection() {
			altered := false
			for _, row := range t.SelectedRows(false) {
				if row.IsOpen() {
					row.SetOpen(false)
					altered = true
				}
			}
			if altered {
				t.SyncToModel()
				t.PruneSelectionOfUndisclosedNodes()
			}
		}
	case KeyRight:
		if t.HasSelection() {
			altered := false
			for _, row := range t.SelectedRows(false) {
				if !row.IsOpen() {
					row.SetOpen(true)
					altered = true
				}
			}
			if altered {
				t.SyncToModel()
			}
		}
	case KeyUp:
		var i int
		if t.HasSelection() {
			i = max(t.FirstSelectedRowIndex()-1, 0)
		} else {
			i = len(t.rowCache) - 1
		}
		if !mod.ShiftDown() {
			t.ClearSelection()
		}
		t.SelectByIndex(i)
		t.ScrollRowCellIntoView(i, 0)
	case KeyDown:
		i := min(t.LastSelectedRowIndex()+1, len(t.rowCache)-1)
		if !mod.ShiftDown() {
			t.ClearSelection()
		}
		t.SelectByIndex(i)
		t.ScrollRowCellIntoView(i, 0)
	case KeyHome:
		if mod.ShiftDown() && t.HasSelection() {
			t.SelectRange(0, t.FirstSelectedRowIndex())
		} else {
			t.ClearSelection()
			t.SelectByIndex(0)
		}
		t.ScrollRowCellIntoView(0, 0)
	case KeyEnd:
		if mod.ShiftDown() && t.HasSelection() {
			t.SelectRange(t.LastSelectedRowIndex(), len(t.rowCache)-1)
		} else {
			t.ClearSelection()
			t.SelectByIndex(len(t.rowCache) - 1)
		}
		t.ScrollRowCellIntoView(len(t.rowCache)-1, 0)
	default:
		return false
	}
	return true
}

// PruneSelectionOfUndisclosedNodes removes any nodes in the selection map that are no longer disclosed from the
// selection map.
func (t *Table[T]) PruneSelectionOfUndisclosedNodes() {
	if !t.selNeedsPrune {
		return
	}
	t.selNeedsPrune = false
	if len(t.selMap) == 0 {
		return
	}
	needsNotify := false
	selMap := make(map[uuid.UUID]bool, len(t.selMap))
	for _, entry := range t.rowCache {
		id := entry.row.UUID()
		if t.selMap[id] {
			selMap[id] = true
		} else {
			needsNotify = true
		}
	}
	t.selMap = selMap
	if needsNotify {
		t.notifyOfSelectionChange()
	}
}

// FirstSelectedRowIndex returns the first selected row index, or -1 if there is no selection.
func (t *Table[T]) FirstSelectedRowIndex() int {
	if len(t.selMap) == 0 {
		return -1
	}
	for i, entry := range t.rowCache {
		if t.selMap[entry.row.UUID()] {
			return i
		}
	}
	return -1
}

// LastSelectedRowIndex returns the last selected row index, or -1 if there is no selection.
func (t *Table[T]) LastSelectedRowIndex() int {
	if len(t.selMap) == 0 {
		return -1
	}
	for i := len(t.rowCache) - 1; i >= 0; i-- {
		if t.selMap[t.rowCache[i].row.UUID()] {
			return i
		}
	}
	return -1
}

// IsRowOrAnyParentSelected returns true if the specified row index or any of its parents are selected.
func (t *Table[T]) IsRowOrAnyParentSelected(index int) bool {
	if index < 0 || index >= len(t.rowCache) {
		return false
	}
	for index >= 0 {
		if t.selMap[t.rowCache[index].row.UUID()] {
			return true
		}
		index = t.rowCache[index].parent
	}
	return false
}

// IsRowSelected returns true if the specified row index is selected.
func (t *Table[T]) IsRowSelected(index int) bool {
	if index < 0 || index >= len(t.rowCache) {
		return false
	}
	return t.selMap[t.rowCache[index].row.UUID()]
}

// SelectedRows returns the currently selected rows. If 'minimal' is true, then children of selected rows that may also
// be selected are not returned, just the topmost row that is selected in any given hierarchy.
func (t *Table[T]) SelectedRows(minimal bool) []T {
	t.PruneSelectionOfUndisclosedNodes()
	if len(t.selMap) == 0 {
		return nil
	}
	rows := make([]T, 0, len(t.selMap))
	for _, entry := range t.rowCache {
		if t.selMap[entry.row.UUID()] && (!minimal || entry.parent == -1 || !t.IsRowOrAnyParentSelected(entry.parent)) {
			rows = append(rows, entry.row)
		}
	}
	return rows
}

// CopySelectionMap returns a copy of the current selection map.
func (t *Table[T]) CopySelectionMap() map[uuid.UUID]bool {
	t.PruneSelectionOfUndisclosedNodes()
	return copySelMap(t.selMap)
}

// SetSelectionMap sets the current selection map.
func (t *Table[T]) SetSelectionMap(selMap map[uuid.UUID]bool) {
	t.selMap = copySelMap(selMap)
	t.selNeedsPrune = true
	t.MarkForRedraw()
	t.notifyOfSelectionChange()
}

func copySelMap(selMap map[uuid.UUID]bool) map[uuid.UUID]bool {
	result := make(map[uuid.UUID]bool, len(selMap))
	for k, v := range selMap {
		result[k] = v
	}
	return result
}

// HasSelection returns true if there is a selection.
func (t *Table[T]) HasSelection() bool {
	t.PruneSelectionOfUndisclosedNodes()
	return len(t.selMap) != 0
}

// SelectionCount returns the number of rows explicitly selected.
func (t *Table[T]) SelectionCount() int {
	t.PruneSelectionOfUndisclosedNodes()
	return len(t.selMap)
}

// ClearSelection clears the selection.
func (t *Table[T]) ClearSelection() {
	if len(t.selMap) == 0 {
		return
	}
	t.selMap = make(map[uuid.UUID]bool)
	t.selNeedsPrune = false
	t.selAnchor = zeroUUID
	t.MarkForRedraw()
	t.notifyOfSelectionChange()
}

// SelectAll selects all rows.
func (t *Table[T]) SelectAll() {
	t.selMap = make(map[uuid.UUID]bool, len(t.rowCache))
	t.selNeedsPrune = false
	t.selAnchor = zeroUUID
	for _, cache := range t.rowCache {
		id := cache.row.UUID()
		t.selMap[id] = true
		if t.selAnchor == zeroUUID {
			t.selAnchor = id
		}
	}
	t.MarkForRedraw()
	t.notifyOfSelectionChange()
}

// SelectByIndex selects the given indexes. The first one will be considered the anchor selection if no existing anchor
// selection exists.
func (t *Table[T]) SelectByIndex(indexes ...int) {
	for _, index := range indexes {
		if index >= 0 && index < len(t.rowCache) {
			id := t.rowCache[index].row.UUID()
			t.selMap[id] = true
			t.selNeedsPrune = true
			if t.selAnchor == zeroUUID {
				t.selAnchor = id
			}
		}
	}
	t.MarkForRedraw()
	t.notifyOfSelectionChange()
}

// SelectRange selects the given range. The start will be considered the anchor selection if no existing anchor
// selection exists.
func (t *Table[T]) SelectRange(start, end int) {
	start = max(start, 0)
	end = min(end, len(t.rowCache)-1)
	if start > end {
		return
	}
	for i := start; i <= end; i++ {
		id := t.rowCache[i].row.UUID()
		t.selMap[id] = true
		t.selNeedsPrune = true
		if t.selAnchor == zeroUUID {
			t.selAnchor = id
		}
	}
	t.MarkForRedraw()
	t.notifyOfSelectionChange()
}

// DeselectByIndex deselects the given indexes.
func (t *Table[T]) DeselectByIndex(indexes ...int) {
	for _, index := range indexes {
		if index >= 0 && index < len(t.rowCache) {
			delete(t.selMap, t.rowCache[index].row.UUID())
		}
	}
	t.MarkForRedraw()
	t.notifyOfSelectionChange()
}

// DeselectRange deselects the given range.
func (t *Table[T]) DeselectRange(start, end int) {
	start = max(start, 0)
	end = min(end, len(t.rowCache)-1)
	if start > end {
		return
	}
	for i := start; i <= end; i++ {
		delete(t.selMap, t.rowCache[i].row.UUID())
	}
	t.MarkForRedraw()
	t.notifyOfSelectionChange()
}

// DiscloseRow ensures the given row can be viewed by opening all parents that lead to it. Returns true if any
// modification was made.
func (t *Table[T]) DiscloseRow(row T, delaySync bool) bool {
	modified := false
	p := row.Parent()
	var zero T
	for p != zero {
		if !p.IsOpen() {
			p.SetOpen(true)
			modified = true
		}
		p = p.Parent()
	}
	if modified {
		if delaySync {
			t.EventuallySyncToModel()
		} else {
			t.SyncToModel()
		}
	}
	return modified
}

// RootRowCount returns the number of top-level rows.
func (t *Table[T]) RootRowCount() int {
	if t.filteredRows != nil {
		return len(t.filteredRows)
	}
	return t.Model.RootRowCount()
}

// RootRows returns the top-level rows. Do not alter the returned list.
func (t *Table[T]) RootRows() []T {
	if t.filteredRows != nil {
		return t.filteredRows
	}
	return t.Model.RootRows()
}

// SetRootRows sets the top-level rows this table will display. This will call SyncToModel() automatically.
func (t *Table[T]) SetRootRows(rows []T) {
	t.filteredRows = nil
	t.Model.SetRootRows(rows)
	t.selMap = make(map[uuid.UUID]bool)
	t.selNeedsPrune = false
	t.selAnchor = zeroUUID
	t.SyncToModel()
}

// SyncToModel causes the table to update its internal caches to reflect the current model.
func (t *Table[T]) SyncToModel() {
	rowCount := 0
	roots := t.RootRows()
	if t.filteredRows != nil {
		rowCount = len(t.filteredRows)
	} else {
		for _, row := range roots {
			rowCount += t.countOpenRowChildrenRecursively(row)
		}
	}
	t.rowCache = make([]tableCache[T], rowCount)
	j := 0
	for _, row := range roots {
		j = t.buildRowCacheEntry(row, -1, j, 0)
	}
	t.selNeedsPrune = true
	_, pref, _ := t.DefaultSizes(Size{})
	rect := t.FrameRect()
	rect.Size = pref
	t.SetFrameRect(rect)
	t.MarkForRedraw()
	t.MarkForLayoutRecursivelyUpward()
}

func (t *Table[T]) countOpenRowChildrenRecursively(row T) int {
	count := 1
	if row.CanHaveChildren() && row.IsOpen() {
		for _, child := range row.Children() {
			count += t.countOpenRowChildrenRecursively(child)
		}
	}
	return count
}

func (t *Table[T]) buildRowCacheEntry(row T, parentIndex, index, depth int) int {
	t.rowCache[index].row = row
	t.rowCache[index].parent = parentIndex
	t.rowCache[index].depth = depth
	t.rowCache[index].height = t.heightForColumns(row, index, depth)
	parentIndex = index
	index++
	if t.filteredRows == nil && row.CanHaveChildren() && row.IsOpen() {
		for _, child := range row.Children() {
			index = t.buildRowCacheEntry(child, parentIndex, index, depth+1)
		}
	}
	return index
}

func (t *Table[T]) heightForColumns(rowData T, row, depth int) float32 {
	var height float32
	for col := range t.Columns {
		w := t.Columns[col].Current
		if w <= 0 {
			continue
		}
		w -= t.Padding.Left + t.Padding.Right
		if t.Columns[col].ID == t.HierarchyColumnID {
			w -= t.Padding.Left + t.HierarchyIndent*float32(depth+1)
		}
		size := t.cellPrefSize(rowData, row, col, w)
		size.Height += t.Padding.Top + t.Padding.Bottom
		if height < size.Height {
			height = size.Height
		}
	}
	return max(xmath.Ceil(height), t.MinimumRowHeight)
}

func (t *Table[T]) cellPrefSize(rowData T, row, col int, widthConstraint float32) geom.Size32 {
	fg, bg, selected, indirectlySelected, focused := t.cellParams(row, col)
	cell := rowData.ColumnCell(row, col, fg, bg, selected, indirectlySelected, focused).AsPanel()
	_, size, _ := cell.Sizes(Size{Width: widthConstraint})
	return size
}

// SizeColumnsToFitWithExcessIn sizes each column to its preferred size, with the exception of the column with the given
// ID, which gets set to any remaining width left over. If the provided column ID doesn't exist, the first column will
// be used instead.
func (t *Table[T]) SizeColumnsToFitWithExcessIn(columnID int) {
	excessColumnIndex := max(t.ColumnIndexForID(columnID), 0)
	current := make([]float32, len(t.Columns))
	for col := range t.Columns {
		current[col] = max(t.Columns[col].Minimum, 0)
		t.Columns[col].Current = 0
	}
	for row, cache := range t.rowCache {
		for col := range t.Columns {
			if col == excessColumnIndex {
				continue
			}
			pref := t.cellPrefSize(cache.row, row, col, 0)
			minimum := t.Columns[col].AutoMinimum
			if minimum > 0 && pref.Width < minimum {
				pref.Width = minimum
			} else {
				maximum := t.Columns[col].AutoMaximum
				if maximum > 0 && pref.Width > maximum {
					pref.Width = maximum
				}
			}
			pref.Width += t.Padding.Left + t.Padding.Right
			if t.Columns[col].ID == t.HierarchyColumnID {
				pref.Width += t.Padding.Left + t.HierarchyIndent*float32(cache.depth+1)
			}
			if current[col] < pref.Width {
				current[col] = pref.Width
			}
		}
	}
	width := t.ContentRect(false).Width
	if t.ShowColumnDivider {
		width -= float32(len(t.Columns) - 1)
	}
	for col := range current {
		if col == excessColumnIndex {
			continue
		}
		t.Columns[col].Current = current[col]
		width -= current[col]
	}
	t.Columns[excessColumnIndex].Current = max(width, t.Columns[excessColumnIndex].Minimum)
	for row, cache := range t.rowCache {
		t.rowCache[row].height = t.heightForColumns(cache.row, row, cache.depth)
	}
}

// SizeColumnsToFit sizes each column to its preferred size. If 'adjust' is true, the Table's FrameRect will be set to
// its preferred size as well.
func (t *Table[T]) SizeColumnsToFit(adjust bool) {
	current := make([]float32, len(t.Columns))
	for col := range t.Columns {
		current[col] = max(t.Columns[col].Minimum, 0)
		t.Columns[col].Current = 0
	}
	for row, cache := range t.rowCache {
		for col := range t.Columns {
			pref := t.cellPrefSize(cache.row, row, col, 0)
			minimum := t.Columns[col].AutoMinimum
			if minimum > 0 && pref.Width < minimum {
				pref.Width = minimum
			} else {
				maximum := t.Columns[col].AutoMaximum
				if maximum > 0 && pref.Width > maximum {
					pref.Width = maximum
				}
			}
			pref.Width += t.Padding.Left + t.Padding.Right
			if t.Columns[col].ID == t.HierarchyColumnID {
				pref.Width += t.Padding.Left + t.HierarchyIndent*float32(cache.depth+1)
			}
			if current[col] < pref.Width {
				current[col] = pref.Width
			}
		}
	}
	for col := range current {
		t.Columns[col].Current = current[col]
	}
	for row, cache := range t.rowCache {
		t.rowCache[row].height = t.heightForColumns(cache.row, row, cache.depth)
	}
	if adjust {
		_, pref, _ := t.DefaultSizes(Size{})
		rect := t.FrameRect()
		rect.Size = pref
		t.SetFrameRect(rect)
	}
}

// SizeColumnToFit sizes the specified column to its preferred size. If 'adjust' is true, the Table's FrameRect will be
// set to its preferred size as well.
func (t *Table[T]) SizeColumnToFit(col int, adjust bool) {
	if col < 0 || col >= len(t.Columns) {
		return
	}
	current := max(t.Columns[col].Minimum, 0)
	t.Columns[col].Current = 0
	for row, cache := range t.rowCache {
		pref := t.cellPrefSize(cache.row, row, col, 0)
		minimum := t.Columns[col].AutoMinimum
		if minimum > 0 && pref.Width < minimum {
			pref.Width = minimum
		} else {
			maximum := t.Columns[col].AutoMaximum
			if maximum > 0 && pref.Width > maximum {
				pref.Width = maximum
			}
		}
		pref.Width += t.Padding.Left + t.Padding.Right
		if t.Columns[col].ID == t.HierarchyColumnID {
			pref.Width += t.Padding.Left + t.HierarchyIndent*float32(cache.depth+1)
		}
		if current < pref.Width {
			current = pref.Width
		}
	}
	t.Columns[col].Current = current
	for row, cache := range t.rowCache {
		t.rowCache[row].height = t.heightForColumns(cache.row, row, cache.depth)
	}
	if adjust {
		_, pref, _ := t.DefaultSizes(Size{})
		rect := t.FrameRect()
		rect.Size = pref
		t.SetFrameRect(rect)
	}
}

// EventuallySizeColumnsToFit sizes each column to its preferred size after a short delay, allowing multiple
// back-to-back calls to this function to only do work once. If 'adjust' is true, the Table's FrameRect will be set to
// its preferred size as well.
func (t *Table[T]) EventuallySizeColumnsToFit(adjust bool) {
	if !t.awaitingSizeColumnsToFit {
		t.awaitingSizeColumnsToFit = true
		InvokeTaskAfter(func() {
			t.SizeColumnsToFit(adjust)
			t.awaitingSizeColumnsToFit = false
		}, 20*time.Millisecond)
	}
}

// EventuallySyncToModel syncs the table to its underlying model after a short delay, allowing multiple back-to-back
// calls to this function to only do work once.
func (t *Table[T]) EventuallySyncToModel() {
	if !t.awaitingSyncToModel {
		t.awaitingSyncToModel = true
		InvokeTaskAfter(func() {
			t.SyncToModel()
			t.awaitingSyncToModel = false
		}, 20*time.Millisecond)
	}
}

// DefaultSizes provides the default sizing.
func (t *Table[T]) DefaultSizes(_ Size) (minSize, prefSize, maxSize Size) {
	for col := range t.Columns {
		prefSize.Width += t.Columns[col].Current
	}
	startRow, endBeforeRow := t.CurrentDrawRowRange()
	for _, cache := range t.rowCache[startRow:endBeforeRow] {
		prefSize.Height += cache.height
	}
	if t.ShowColumnDivider {
		prefSize.Width += float32(len(t.Columns) - 1)
	}
	if t.ShowRowDivider {
		prefSize.Height += float32((endBeforeRow - startRow) - 1)
	}
	if border := t.Border(); border != nil {
		prefSize.AddInsets(border.Insets())
	}
	prefSize.GrowToInteger()
	return prefSize, prefSize, prefSize
}

// RowFromIndex returns the row data for the given index.
func (t *Table[T]) RowFromIndex(index int) T {
	if index < 0 || index >= len(t.rowCache) {
		var zero T
		return zero
	}
	return t.rowCache[index].row
}

// RowToIndex returns the row's index within the displayed data, or -1 if it isn't currently in the disclosed rows.
func (t *Table[T]) RowToIndex(rowData T) int {
	id := rowData.UUID()
	for row, data := range t.rowCache {
		if data.row.UUID() == id {
			return row
		}
	}
	return -1
}

// LastRowIndex returns the index of the last row. Will be -1 if there are no rows.
func (t *Table[T]) LastRowIndex() int {
	return len(t.rowCache) - 1
}

// ScrollRowIntoView scrolls the row at the given index into view.
func (t *Table[T]) ScrollRowIntoView(row int) {
	if frame := t.RowFrame(row); !frame.IsEmpty() {
		t.ScrollRectIntoView(frame)
	}
}

// ScrollRowCellIntoView scrolls the cell from the row and column at the given indexes into view.
func (t *Table[T]) ScrollRowCellIntoView(row, col int) {
	if frame := t.CellFrame(row, col); !frame.IsEmpty() {
		t.ScrollRectIntoView(frame)
	}
}

// IsFiltered returns true if a filter is currently applied. When a filter is applied, no hierarchy is display and no
// modifications to the row data should be performed.
func (t *Table[T]) IsFiltered() bool {
	return t.filteredRows != nil
}

// ApplyFilter applies a filter to the data. When a non-nil filter is applied, all rows (recursively) are passed through
// the filter. Only those that the filter returns false for will be visible in the table. When a filter is applied, no
// hierarchy is display and no modifications to the row data should be performed.
func (t *Table[T]) ApplyFilter(filter func(row T) bool) {
	if filter == nil {
		if t.filteredRows == nil {
			return
		}
		t.filteredRows = nil
	} else {
		t.filteredRows = make([]T, 0)
		for _, row := range t.Model.RootRows() {
			t.applyFilter(row, filter)
		}
	}
	t.SyncToModel()
	if t.header != nil && t.header.HasSort() {
		t.header.ApplySort()
	}
}

func (t *Table[T]) applyFilter(row T, filter func(row T) bool) {
	if !filter(row) {
		t.filteredRows = append(t.filteredRows, row)
	}
	if row.CanHaveChildren() {
		for _, child := range row.Children() {
			t.applyFilter(child, filter)
		}
	}
}

// InstallDragSupport installs default drag support into a table. This will chain a function to any existing
// MouseDragCallback.
func (t *Table[T]) InstallDragSupport(svg *SVG, dragKey, singularName, pluralName string) {
	orig := t.MouseDragCallback
	t.MouseDragCallback = func(where Point, button int, mod Modifiers) bool {
		if orig != nil && orig(where, button, mod) {
			return true
		}
		if button == ButtonLeft && t.HasSelection() && t.IsDragGesture(where) {
			data := &TableDragData[T]{
				Table: t,
				Rows:  t.SelectedRows(true),
			}
			drawable := NewTableDragDrawable(data, svg, singularName, pluralName)
			size := drawable.LogicalSize()
			t.StartDataDrag(&DragData{
				Data:     map[string]any{dragKey: data},
				Drawable: drawable,
				Ink:      t.OnBackgroundInk,
				Offset:   Point{X: 0, Y: -size.Height / 2},
			})
		}
		return false
	}
}

// InstallDropSupport installs default drop support into a table. This will replace any existing DataDragOverCallback,
// DataDragExitCallback, and DataDragDropCallback functions. It will also chain a function to any existing
// DrawOverCallback. The shouldMoveDataCallback is called when a drop is about to occur to determine if the data should
// be moved (i.e. removed from the source) or copied to the destination. The willDropCallback is called before the
// actual data changes are made, giving an opportunity to start an undo event, which should be returned. The
// didDropCallback is called after data changes are made and is passed the undo event (if any) returned by the
// willDropCallback, so that the undo event can be completed and posted.
func InstallDropSupport[T TableRowConstraint[T], U any](t *Table[T], dragKey string, shouldMoveDataCallback func(from, to *Table[T]) bool, willDropCallback func(from, to *Table[T], move bool) *UndoEdit[U], didDropCallback func(undo *UndoEdit[U], from, to *Table[T], move bool)) *TableDrop[T, U] {
	drop := &TableDrop[T, U]{
		Table:                  t,
		DragKey:                dragKey,
		originalDrawOver:       t.DrawOverCallback,
		shouldMoveDataCallback: shouldMoveDataCallback,
		willDropCallback:       willDropCallback,
		didDropCallback:        didDropCallback,
	}
	t.DataDragOverCallback = drop.DataDragOverCallback
	t.DataDragExitCallback = drop.DataDragExitCallback
	t.DataDragDropCallback = drop.DataDragDropCallback
	t.DrawOverCallback = drop.DrawOverCallback
	return drop
}

// CountTableRows returns the number of table rows, including all descendants, whether open or not.
func CountTableRows[T TableRowConstraint[T]](rows []T) int {
	count := len(rows)
	for _, row := range rows {
		if row.CanHaveChildren() {
			count += CountTableRows(row.Children())
		}
	}
	return count
}

// RowContainsRow returns true if 'descendant' is in fact a descendant of 'ancestor'.
func RowContainsRow[T TableRowConstraint[T]](ancestor, descendant T) bool {
	var zero T
	for descendant != zero && descendant != ancestor {
		descendant = descendant.Parent()
	}
	return descendant == ancestor
}
