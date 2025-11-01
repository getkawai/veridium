package document_test

import (
	"testing"

	"github.com/kawai-network/veridium/gooxml/document"
)

func TestToMarkdown(t *testing.T) {
	doc := document.New()

	// Add a heading
	para := doc.AddParagraph()
	para.SetStyle("Heading1")
	run := para.AddRun()
	run.AddText("Test Heading")

	// Add a paragraph with bold, italic, strikethrough, and code
	para = doc.AddParagraph()
	run = para.AddRun()
	run.AddText("Normal text ")
	run = para.AddRun()
	run.Properties().SetBold(true)
	run.AddText("bold text")
	run = para.AddRun()
	run.AddText(" ")
	run = para.AddRun()
	run.Properties().SetItalic(true)
	run.AddText("italic text")
	run = para.AddRun()
	run.AddText(" ")
	run = para.AddRun()
	run.Properties().SetStrikeThrough(true)
	run.AddText("strikethrough")
	run = para.AddRun()
	run.AddText(" ")
	run = para.AddRun()
	run.Properties().SetFontFamily("Courier New")
	run.AddText("monospace")

	// Add a blockquote (indented paragraph)
	para = doc.AddParagraph()
	para.Properties().SetStartIndent(720) // 0.5 inches
	run = para.AddRun()
	run.AddText("This is a blockquote")

	// Add a table
	tbl := doc.AddTable()
	row := tbl.AddRow()
	cell := row.AddCell()
	para = cell.AddParagraph()
	para.AddRun().AddText("Header 1")
	cell = row.AddCell()
	para = cell.AddParagraph()
	para.AddRun().AddText("Header 2")

	row = tbl.AddRow()
	cell = row.AddCell()
	para = cell.AddParagraph()
	para.AddRun().AddText("Data 1")
	cell = row.AddCell()
	para = cell.AddParagraph()
	para.AddRun().AddText("Data 2")

	md := doc.ToMarkdown()
	expected := `# Test Heading

Normal text **bold text** *italic text* ~~strikethrough~~ monospace

> This is a blockquote

| Header 1 | Header 2 |
| --- | --- |
| Data 1 | Data 2 |

`
	if md != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, md)
	}
}
