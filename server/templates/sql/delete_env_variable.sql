DELETE FROM environment_variables
WHERE application = {{ .AppId | __sql_arg__ }} AND key = {{ .Key | __sql_arg__ }}
