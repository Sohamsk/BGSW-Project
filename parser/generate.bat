@echo off
:: Set the maximum heap size and specify the classpath
set JAVA_OPTS=-Xmx500M
set CLASSPATH=antlr.jar;%CLASSPATH%

:: Run the ANTLR Tool
java %JAVA_OPTS% -cp "%CLASSPATH%" org.antlr.v4.Tool -Dlanguage=Go -visitor -package parser *.g4
