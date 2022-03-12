package main

var stylesheet = `
body { 
	margin: 0; 
}
td {
	vertical-align: top;
}
#top-bar .tab-bar-divider {
	border-bottom: 1px solid gainsboro;
	position: relative;
	top: -3px;
	z-index: -1;
}
#top-tab .app-title {
	padding: 0 16px 0 16px
}
.mdc-layout-grid__inner {
	margin-top: 5px;
}
.mdc-data-table__header-cell {
	font-weight: bold;
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
`
