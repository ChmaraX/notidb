package notion

import (
	"strings"

	"github.com/ChmaraX/notidb/internal/utils"
	"github.com/jomei/notionapi"
)

func GetSupportedPropTypes() []notionapi.PropertyType {
	return []notionapi.PropertyType{
		notionapi.PropertyTypeTitle,
		notionapi.PropertyTypeRichText,
		notionapi.PropertyTypeNumber,
		notionapi.PropertyTypeSelect,
		notionapi.PropertyTypeMultiSelect,
		notionapi.PropertyTypeDate,
		notionapi.PropertyTypeCheckbox,
		notionapi.PropertyTypeEmail,
		notionapi.PropertyTypePhoneNumber,
	}
}

func CreateContentBlock(content string) notionapi.Block {
	return notionapi.ParagraphBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: "block",
			Type:   "paragraph",
		},
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{
					Type: "text",
					Text: &notionapi.Text{
						Content: content,
					},
				},
			},
		},
	}
}

func CreateTitleProperty(title string) notionapi.TitleProperty {
	return notionapi.TitleProperty{Title: []notionapi.RichText{
		{
			Type: "text",
			Text: &notionapi.Text{
				Content: title,
			},
		},
	},
	}
}

func CreateRichTextProperty(content string) notionapi.RichTextProperty {
	return notionapi.RichTextProperty{
		RichText: []notionapi.RichText{
			{
				Type: "text",
				Text: &notionapi.Text{
					Content: content,
				},
			},
		},
	}
}

func CreateNumberProperty(number float64) notionapi.NumberProperty {
	return notionapi.NumberProperty{Number: number}
}

func CreateSelectProperty(option string) notionapi.SelectProperty {
	return notionapi.SelectProperty{Select: notionapi.Option{Name: option}}
}

func CreateMultiSelectProperty(options []string) notionapi.MultiSelectProperty {
	var opts []notionapi.Option
	for _, option := range options {
		if strings.TrimSpace(option) != "" {
			opts = append(opts, notionapi.Option{Name: option})
		}
	}
	return notionapi.MultiSelectProperty{MultiSelect: opts}
}

func CreateDateProperty(date string) (notionapi.DateProperty, error) {
	// if dateString doesn't contain time - default to 12:00 AM
	if !strings.Contains(date, ":") {
		date = date + " 00:00"
	}
	dateTime, err := utils.ParseDateInLocation(date, "02/01/2006 15:04")
	if err != nil {
		return notionapi.DateProperty{}, err
	}
	start := notionapi.Date(dateTime)
	return notionapi.DateProperty{Date: &notionapi.DateObject{Start: &start}}, nil
}

func CreateCheckboxProperty(checked bool) notionapi.CheckboxProperty {
	return notionapi.CheckboxProperty{Checkbox: checked}
}

func CreateEmailProperty(email string) notionapi.EmailProperty {
	return notionapi.EmailProperty{Email: email}
}

func CreatePhoneNumberProperty(phoneNumber string) notionapi.PhoneNumberProperty {
	return notionapi.PhoneNumberProperty{PhoneNumber: phoneNumber}
}
