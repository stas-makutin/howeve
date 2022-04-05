package main

var stylesheet = `
body { 
	margin: 0; 
}
#top-tab .tab-bar-divider {
	border-bottom: 1px solid gainsboro;
	position: relative;
	top: -3px;
	z-index: -1;
}
#top-tab .app-title {
	padding: 0 16px 0 16px
}
.data-table-cell--top {
	vertical-align: top;
}
.align-center {
	text-align: center;
}
.mdc-data-table__header-cell {
	font-weight: bold;
}
.adjacent-margins {
	margin: 0.5em 0.5em 0.5em 0;
}
.view {
	position: relative;
}
.view-loading {
	position: absolute;
	top: 0;
	left: 0;
	height: 100%;
	width: 100%;
	zIndex: 100;
	vertical-align: middle;
	opacity: 70%;
}
.view-loading__progress {
	margin: 10% 50%;
	opacity: 100%;
}
#sv-add-service-dialog---content {
	min-height: 15em;
}
#sv-add-service-protocols {
	min-width: 18em;
	margin-top: 0.5em;
	margin-right: 0.5em;
}
#sv-add-service-transports {
	min-width: 18em;
	margin-top: 0.5em;
	margin-right: 0.5em;
}
#sv-add-service-alias {
	min-width: 18em;
	margin-top: 0.5em;
}
#sv-add-service-entry {
	min-width: 16em;
	width: 100%;
	margin-top: 0.5em;
}
#sv-add-service-add-param {
	margin-top: 0.5em;
	margin-bottom: 0.5em;
}
.sv-add-service-param-name {
	min-width: 15%;
	width: 25%;
	margin-top: 0.5em;
	margin-right: 0.5em;
}
.sv-add-service-param-value {
	min-width: 15%;
	width: 55%;
	margin-top: 0.5em;
}
.sv-add-service-param-delete {
	position: relative;
	top: 0.2em;
}
`
