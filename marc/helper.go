package marc

type TagsToReadableMap map[string]string

var pTagsToReadableMap *TagsToReadableMap = new(TagsToReadableMap)

func setTagsToReadable(fieldTag string, indicator string, sTag string, content string) {
	(*pTagsToReadableMap)[fieldTag+sTag] = content
}

func getTagsToReadable(fieldTag string, indicator string, sTag string) string {
	return (*pTagsToReadableMap)[fieldTag+sTag]
}
