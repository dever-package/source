package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shemic/dever/server"
	"github.com/shemic/dever/util"

	frontaction "github.com/dever-package/front/service/action"
	frontrecord "github.com/dever-package/front/service/record"
	sourcemodel "github.com/dever-package/source/model"
)

type SourceHook struct{}
type SourceOptionService struct{}

const (
	channelModelName = "source.NewChannelModel"
	cateModelName    = "source.NewCateModel"
	originModelName  = "source.NewOriginModel"
	statusModelName  = "source.NewStatusModel"
)

func (SourceHook) ProviderBeforeSaveSource(c *server.Context, params []any) any {
	record := cloneSourceRecord(params)
	if len(record) == 0 {
		return record
	}
	if titleValue, ok := record["title"]; ok {
		record["title"] = strings.TrimSpace(util.ToString(titleValue))
		if record["title"] == "" {
			panicField("form.title", "标题不能为空。")
		}
	}
	if isPartialRecord(record) {
		normalizeSourceImages(record)
		normalizePresentSort(record)
		normalizePresentEnabledStatus(record)
		normalizePresentSourceStatusID(c, record)
		return record
	}

	if record["title"] == "" {
		panicField("form.title", "标题不能为空。")
	}

	cateID := util.ToUint64(record["cate_id"])
	if cateID == 0 {
		panicField("form.cate_id", "资源必须选择一个分类。")
	}

	cateRows := loadActiveRows(c.Context(), cateModelName, []uint64{cateID})
	if len(cateRows) != 1 {
		panicField("form.cate_id", "资源选择的分类不存在或已停用。")
	}
	channelID := util.ToUint64(cateRows[0]["channel_id"])
	if channelID == 0 {
		panicField("form.cate_id", "分类缺少所属频道。")
	}
	selectedChannelID := util.ToUint64(record["channel_id"])
	if selectedChannelID > 0 && selectedChannelID != channelID {
		panicField("form.cate_id", "分类不属于所选频道。")
	}

	channelRows := loadActiveRows(c.Context(), channelModelName, []uint64{channelID})
	if len(channelRows) != 1 {
		panicField("form.cate_id", "分类所属频道不存在或已停用。")
	}

	record["channel_id"] = channelID
	record["cate_id"] = cateID
	normalizeOrigin(c, record)
	normalizeSourceImages(record)
	normalizeFullSort(record)
	normalizeFullEnabledStatus(record)
	normalizeSourceStatusID(c, record)
	return record
}

func (SourceHook) ProviderBeforeSaveChannel(_ *server.Context, params []any) any {
	record := cloneSourceRecord(params)
	if nameValue, ok := record["name"]; ok {
		record["name"] = strings.TrimSpace(util.ToString(nameValue))
		if record["name"] == "" {
			panicField("form.name", "频道名称不能为空。")
		}
	}
	if !isPartialRecord(record) && strings.TrimSpace(util.ToString(record["name"])) == "" {
		panicField("form.name", "频道名称不能为空。")
	}
	if isPartialRecord(record) {
		normalizePresentSortAndEnabledStatus(record)
	} else {
		normalizeFullSortAndEnabledStatus(record)
	}
	return record
}

func (SourceHook) ProviderBeforeSaveOrigin(_ *server.Context, params []any) any {
	record := cloneSourceRecord(params)
	if nameValue, ok := record["name"]; ok {
		record["name"] = strings.TrimSpace(util.ToString(nameValue))
		if record["name"] == "" {
			panicField("form.name", "来源名称不能为空。")
		}
	}
	if !isPartialRecord(record) && strings.TrimSpace(util.ToString(record["name"])) == "" {
		panicField("form.name", "来源名称不能为空。")
	}
	if isPartialRecord(record) {
		normalizePresentSortAndEnabledStatus(record)
	} else {
		normalizeFullSortAndEnabledStatus(record)
	}
	return record
}

func (SourceHook) ProviderBeforeSaveStatus(_ *server.Context, params []any) any {
	record := cloneSourceRecord(params)
	if nameValue, ok := record["name"]; ok {
		record["name"] = strings.TrimSpace(util.ToString(nameValue))
		if record["name"] == "" {
			panicField("form.name", "状态名称不能为空。")
		}
	}
	if !isPartialRecord(record) && strings.TrimSpace(util.ToString(record["name"])) == "" {
		panicField("form.name", "状态名称不能为空。")
	}
	if colorValue, ok := record["color"]; ok {
		record["color"] = strings.TrimSpace(util.ToString(colorValue))
	}
	if isPartialRecord(record) {
		normalizePresentSortAndEnabledStatus(record)
	} else {
		normalizeFullSortAndEnabledStatus(record)
	}
	return record
}

func (SourceHook) ProviderAfterSaveStatus(_ *server.Context, params []any) any {
	return cloneSourceRecord(params)
}

func (SourceHook) ProviderBeforeSaveCate(c *server.Context, params []any) any {
	record := cloneSourceRecord(params)
	cateID := util.ToUint64(record["id"])
	if nameValue, ok := record["name"]; ok {
		record["name"] = strings.TrimSpace(util.ToString(nameValue))
		if record["name"] == "" {
			panicField("form.name", "分类名称不能为空。")
		}
	}
	if !isPartialRecord(record) && strings.TrimSpace(util.ToString(record["name"])) == "" {
		panicField("form.name", "分类名称不能为空。")
	}

	if isPartialRecord(record) {
		normalizePresentSortAndEnabledStatus(record)
		return record
	}

	channelID := util.ToUint64(record["channel_id"])
	if channelID == 0 {
		panicField("form.channel_id", "分类必须选择所属频道。")
	}
	if len(loadActiveRows(c.Context(), channelModelName, []uint64{channelID})) != 1 {
		panicField("form.channel_id", "分类选择的频道不存在或已停用。")
	}

	cateModel := frontrecord.Resolve(cateModelName)
	if cateModel == nil {
		panic("分类模型未注册")
	}
	if cateID > 0 {
		if len(cateModel.FindMap(c.Context(), map[string]any{"id": cateID})) == 0 {
			panic("分类不存在")
		}
	}

	normalizeFullSortAndEnabledStatus(record)
	return record
}

func (SourceHook) ProviderBeforeDeleteChannel(c *server.Context, params []any) any {
	payload := cloneSourceRecord(params)
	channelID := util.ToUint64(payload["id"])
	if channelID == 0 {
		panic("频道不存在")
	}

	cateModel := frontrecord.Resolve(cateModelName)
	if cateModel != nil && cateModel.Count(c.Context(), map[string]any{"channel_id": channelID}) > 0 {
		panic("当前频道下仍有分类，请先处理分类后再删除")
	}
	sourceModel := frontrecord.Resolve("source.NewSourceModel")
	if sourceModel != nil && sourceModel.Count(c.Context(), map[string]any{"channel_id": channelID}) > 0 {
		panic("当前频道下仍有关联资源，请先处理资源后再删除")
	}
	return map[string]any{"id": channelID}
}

func (SourceHook) ProviderBeforeDeleteCate(c *server.Context, params []any) any {
	payload := cloneSourceRecord(params)
	cateID := util.ToUint64(payload["id"])
	if cateID == 0 {
		panic("分类不存在")
	}

	sourceModel := frontrecord.Resolve("source.NewSourceModel")
	if sourceModel != nil && sourceModel.Count(c.Context(), map[string]any{"cate_id": cateID}) > 0 {
		panic("当前分类下仍有关联资源，请先处理资源后再删除")
	}
	return map[string]any{"id": cateID}
}

func (SourceHook) ProviderBeforeDeleteOrigin(c *server.Context, params []any) any {
	payload := cloneSourceRecord(params)
	originID := util.ToUint64(payload["id"])
	if originID == 0 {
		panic("来源不存在")
	}
	if originID == 1 {
		panic("默认来源不能删除。")
	}
	sourceModel := frontrecord.Resolve("source.NewSourceModel")
	if sourceModel != nil && sourceModel.Count(c.Context(), map[string]any{"origin_id": originID}) > 0 {
		panic("当前来源已有资源使用，不能删除。")
	}
	return map[string]any{"id": originID}
}

func (SourceHook) ProviderBeforeDeleteStatus(c *server.Context, params []any) any {
	payload := cloneSourceRecord(params)
	statusID := util.ToUint64(payload["id"])
	if statusID == 0 {
		panic("资源状态不存在")
	}
	sourceModel := frontrecord.Resolve("source.NewSourceModel")
	if sourceModel != nil && sourceModel.Count(c.Context(), map[string]any{"status_id": statusID}) > 0 {
		panic("当前资源状态已有资源使用，不能删除。")
	}
	return map[string]any{"id": statusID}
}

func (SourceOptionService) ProviderLoadCates(c *server.Context, params []any) any {
	payload := map[string]any{}
	if len(params) > 0 {
		payload, _ = params[0].(map[string]any)
	}

	channelIDs := normalizeIDList(firstPresent(payload, "channel_id", "parent_id"))
	parentField := strings.TrimSpace(util.ToString(payload["parent_field"]))
	if parentField == "channel_id" && len(channelIDs) == 0 {
		return []map[string]any{}
	}
	filters := map[string]any{"status": 1}
	if len(channelIDs) > 0 {
		filters["channel_id"] = uint64sToAny(channelIDs)
	}

	cateModel := frontrecord.Resolve(cateModelName)
	if cateModel == nil {
		return []map[string]any{}
	}

	rows := cateModel.SelectMap(c.Context(), filters, map[string]any{"order": "channel_id asc,sort asc,id asc"})
	if _, ok := payload["parent_id"]; ok {
		return buildCateCascaderRows(rows)
	}
	return buildCateOptionRows(c.Context(), rows)
}

func (SourceOptionService) ProviderLoadChannels(c *server.Context, _ []any) any {
	channelModel := frontrecord.Resolve(channelModelName)
	if channelModel == nil {
		return []map[string]any{}
	}

	rows := channelModel.SelectMap(c.Context(), map[string]any{"status": 1}, map[string]any{
		"field": "main.id, main.name, main.status, main.sort",
		"order": "main.sort asc, main.id asc",
	})
	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(util.ToString(row["name"]))
		result = append(result, map[string]any{
			"id":     util.ToUint64(row["id"]),
			"value":  name,
			"label":  name,
			"name":   name,
			"status": util.ToIntDefault(row["status"], 0),
			"sort":   util.ToIntDefault(row["sort"], 0),
		})
	}
	return result
}

func buildCateCascaderRows(rows []map[string]any) []map[string]any {
	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(util.ToString(row["name"]))
		result = append(result, map[string]any{
			"id":         util.ToUint64(row["id"]),
			"value":      name,
			"label":      name,
			"name":       name,
			"channel_id": util.ToUint64(row["channel_id"]),
		})
	}
	return result
}

func cloneSourceRecord(params []any) map[string]any {
	if len(params) == 0 {
		return map[string]any{}
	}
	record, _ := params[0].(map[string]any)
	if record == nil {
		return map[string]any{}
	}
	return util.CloneMap(record)
}

func isPartialRecord(record map[string]any) bool {
	return util.ToBool(record["_partial"])
}

func normalizePresentSortAndEnabledStatus(record map[string]any) {
	if sortValue, ok := record["sort"]; ok && util.ToIntDefault(sortValue, 0) <= 0 {
		record["sort"] = 100
	}
	if statusValue, ok := record["status"]; ok {
		record["status"] = normalizeEnabledStatus(statusValue)
	}
}

func normalizeFullSortAndEnabledStatus(record map[string]any) {
	normalizePresentSortAndEnabledStatus(record)
	if _, ok := record["sort"]; !ok && util.ToUint64(record["id"]) == 0 {
		record["sort"] = 100
	}
	if _, ok := record["status"]; !ok && util.ToUint64(record["id"]) == 0 {
		record["status"] = 1
	}
}

func normalizePresentSort(record map[string]any) {
	if sortValue, ok := record["sort"]; ok && util.ToIntDefault(sortValue, 0) <= 0 {
		record["sort"] = 100
	}
}

func normalizeFullSort(record map[string]any) {
	normalizePresentSort(record)
	if _, ok := record["sort"]; !ok && util.ToUint64(record["id"]) == 0 {
		record["sort"] = 100
	}
}

func normalizePresentEnabledStatus(record map[string]any) {
	if statusValue, ok := record["status"]; ok {
		record["status"] = normalizeEnabledStatus(statusValue)
	}
}

func normalizeFullEnabledStatus(record map[string]any) {
	normalizePresentEnabledStatus(record)
	if _, ok := record["status"]; !ok && util.ToUint64(record["id"]) == 0 {
		record["status"] = 1
	}
}

func normalizeOrigin(c *server.Context, record map[string]any) {
	_, hasOrigin := record["origin_id"]
	sourceID := util.ToUint64(record["id"])
	originID := util.ToUint64(record["origin_id"])
	if originID == 0 && sourceID != 0 && !hasOrigin {
		return
	}
	if originID == 0 && sourceID != 0 && hasOrigin {
		panicField("form.origin_id", "资源必须选择一个来源。")
	}
	if originID == 0 {
		originID = sourcemodel.DefaultOriginID
	}
	if len(loadActiveRows(c.Context(), originModelName, []uint64{originID})) != 1 {
		panicField("form.origin_id", "资源选择的来源不存在或已停用。")
	}
	record["origin_id"] = originID
}

func normalizePresentSourceStatusID(c *server.Context, record map[string]any) {
	if _, ok := record["status_id"]; !ok {
		return
	}
	normalizeSourceStatusID(c, record)
}

func normalizeSourceStatusID(c *server.Context, record map[string]any) {
	_, hasStatusID := record["status_id"]
	sourceID := util.ToUint64(record["id"])
	statusID := util.ToUint64(record["status_id"])
	if statusID == 0 && sourceID != 0 && !hasStatusID {
		return
	}
	if statusID == 0 && sourceID != 0 && hasStatusID {
		panicField("form.status_id", "资源必须选择一个状态。")
	}
	if statusID == 0 {
		statusID = sourcemodel.DefaultStatusID
	}
	if len(loadActiveRows(c.Context(), statusModelName, []uint64{statusID})) != 1 {
		panicField("form.status_id", "资源选择的状态不存在或已停用。")
	}
	record["status_id"] = statusID
}

func normalizeSourceImages(record map[string]any) {
	if imageValue, ok := record["images"]; ok {
		record["images"] = normalizeUploadJSON(imageValue)
	}
}

func normalizeUploadJSON(value any) string {
	switch current := value.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(current)
	case []map[string]any, []any, map[string]any:
		encoded, err := json.Marshal(current)
		if err == nil {
			return string(encoded)
		}
	}
	return strings.TrimSpace(util.ToString(value))
}

func firstPresent(record map[string]any, keys ...string) any {
	for _, key := range keys {
		if value, ok := record[key]; ok {
			return value
		}
	}
	return nil
}

func buildCateOptionRows(ctx context.Context, rows []map[string]any) []map[string]any {
	channelModel := frontrecord.Resolve(channelModelName)
	channelNames := map[uint64]string{}
	if channelModel != nil {
		channelIDs := make([]any, 0, len(rows))
		seen := map[uint64]struct{}{}
		for _, row := range rows {
			channelID := util.ToUint64(row["channel_id"])
			if channelID == 0 {
				continue
			}
			if _, ok := seen[channelID]; ok {
				continue
			}
			seen[channelID] = struct{}{}
			channelIDs = append(channelIDs, channelID)
		}
		for _, row := range channelModel.SelectMap(ctx, map[string]any{"id": channelIDs}) {
			channelNames[util.ToUint64(row["id"])] = strings.TrimSpace(util.ToString(row["name"]))
		}
	}

	result := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		channelID := util.ToUint64(row["channel_id"])
		name := strings.TrimSpace(util.ToString(row["name"]))
		if channelName := channelNames[channelID]; channelName != "" {
			name = fmt.Sprintf("%s / %s", channelName, name)
		}
		result = append(result, map[string]any{
			"id":         row["id"],
			"value":      name,
			"label":      name,
			"name":       name,
			"channel_id": channelID,
		})
	}
	return result
}

func loadActiveRows(ctx context.Context, modelName string, ids []uint64) []map[string]any {
	model := frontrecord.Resolve(modelName)
	if model == nil || len(ids) == 0 {
		return nil
	}
	idValues := make([]any, 0, len(ids))
	for _, id := range ids {
		idValues = append(idValues, id)
	}
	return model.SelectMap(ctx, map[string]any{
		"id":     idValues,
		"status": 1,
	})
}

func normalizeIDList(value any) []uint64 {
	switch current := value.(type) {
	case []uint64:
		return util.UniqueUint64s(current)
	case []int:
		result := make([]uint64, 0, len(current))
		for _, item := range current {
			if item > 0 {
				result = append(result, uint64(item))
			}
		}
		return util.UniqueUint64s(result)
	case []any:
		result := make([]uint64, 0, len(current))
		for _, item := range current {
			if id := util.ToUint64(item); id > 0 {
				result = append(result, id)
			}
		}
		return util.UniqueUint64s(result)
	case string:
		parts := strings.Split(current, ",")
		result := make([]uint64, 0, len(parts))
		for _, part := range parts {
			if id := util.ToUint64(strings.TrimSpace(part)); id > 0 {
				result = append(result, id)
			}
		}
		return util.UniqueUint64s(result)
	default:
		if id := util.ToUint64(value); id > 0 {
			return []uint64{id}
		}
		return nil
	}
}

func uint64sToAny(ids []uint64) []any {
	result := make([]any, 0, len(ids))
	for _, id := range ids {
		result = append(result, id)
	}
	return result
}

func normalizeEnabledStatus(value any) int16 {
	status := int16(util.ToIntDefault(value, 1))
	if status != 1 && status != 2 {
		return 1
	}
	return status
}

func panicField(field string, message string) {
	panic(frontaction.NewFieldError(field, message))
}
