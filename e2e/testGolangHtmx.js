'use strict';

var testSuite = require('./test.js');
var argv = require('optimist').default('laxMode', false).default('browser', 'chrome').argv;
var frameworkPathLookup = require('./framework-path-lookup');
var rootUrl = 'http://localhost:8080/';

testSuite.todoMVCTest("Golang+htmx", rootUrl, argv.speedMode, argv.laxMode, argv.browser);
