package parser


import (
    "html/template"
    "io"
    "os"
    "fmt"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
)


// For ease of unit testing.
var (
    parseFile           = template.ParseFiles
    openFile            = os.Open
    createFile          = os.Create
    ioReadAll           = io.ReadAll
    yamlUnmarshal       = yaml.Unmarshal
    executeTemplateFile = func(templateFile *template.Template, wr io.Writer, data any) error {
        return templateFile.Execute(wr, data)
    }
)


// valuesFromYamlFile extracts values from yaml file.
func valuesFromYamlFile(dataFile string) (map[string]interface{}, error) {
    data, err := openFile(dataFile)
    if err != nil {
        return nil, errors.Wrap(err, "opening data file")
    }
    defer data.Close()
    s, err := ioReadAll(data)
    if err != nil {
        return nil, errors.Wrap(err, "reading data file")
    }
    var values map[string]interface{}
    err = yamlUnmarshal(s, &values)
    if err != nil {
        return nil, errors.Wrap(err, "unmarshalling yaml file")
    }
    return values, nil
}


// Parse replaces values present in the template file
// with values defined in the data file, saving the result
// as an output file.
func Parse(templateFile, dataFile, outputFile string, targetDir string) error {
    tmpl, err := parseFile(templateFile)
    if err != nil {
        return errors.Wrap(err, "parsing template file")
    }
    values, err := valuesFromYamlFile(dataFile)
    if err != nil {
        return err
    }
    _ = os.Mkdir(targetDir, os.ModePerm)
    targetPath := fmt.Sprintf("%s/%s",targetDir,outputFile)
    output, err := createFile(targetPath)
    if err != nil {
        return errors.Wrap(err, "creating output file")
    }
    defer output.Close()
    err = executeTemplateFile(tmpl, output, values)
    if err != nil {
        return errors.Wrap(err, "executing template file")
    }
    return nil
}

func ParseValues(templateFile string, values map[string]string, outputFile string, targetDir string) error {
//     if values == nil {
//         return errors.Wrap(err, "values not defined")
//     }
    tmpl, err := parseFile(templateFile)
    if err != nil {
        return errors.Wrap(err, "parsing template file")
    }
    _ = os.Mkdir(targetDir, os.ModePerm)
    targetPath := fmt.Sprintf("%s/%s",targetDir,outputFile)
    output, err := createFile(targetPath)
    if err != nil {
        return errors.Wrap(err, "creating output file")
    }
    defer output.Close()
    err = executeTemplateFile(tmpl, output, values)
    if err != nil {
        return errors.Wrap(err, "executing template file")
    }
    return nil
}