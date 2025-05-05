DELETE FROM environment_variables
WHERE application = {{ .AppId }} AND key = {{ .Key }}
