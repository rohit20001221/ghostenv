DELETE FROM applications 
WHERE app_name = {{ .appName | __sql_arg__ }};