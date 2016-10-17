var React = require('react');
var ReactCSSTransitionGroup = require('react-addons-css-transition-group');
var ReactTransitionGroup = require('react-addons-transition-group');
var ReactDOM = require('react-dom');
var ClassNames = require('classnames');

var ExpanderWrapper = React.createClass({
    getInitialState: function() {
        return {
            expanded: false,
        };
    },

    componentWillUnmount: function() {
        this.props.unmounted();
    },

    render: function() {
        return (
            <div className={ClassNames({
                    "expander": true,
                    "xyz": true,
                    "expanded": this.state.expanded })}>
                {this.props.content}
            </div>
        );
    }
});

var Grid = React.createClass({
    getDefaultProps: function() {
        return {
            cells: [],
            cellWidth: 110,
            cellHeight: 110,
            paddingLeft: 0,
            paddingRight: 0,
            paddingTop: 0,
            paddingBottom: 12,
            spacingH: 0,
            spacingV: 0
        };
    },

    getInitialState: function() {
        return {
            cellIndices: { x:-1, y:-1 },
            expanderIndex: -1,
            selectedIndex: -1,
            expanded: false,
            cellWidth: this.props.cellWidth,
            cellHeight: this.props.cellHeight
        }
    },

    calcCellDimensions: function() {
        var $this = $(ReactDOM.findDOMNode(this));
        var gridWidthNoPadding = $this.width();
        var cellsPerRow = Math.floor(gridWidthNoPadding / this.props.cellWidth);

        var fittedWidth = Math.floor(this.props.cellWidth + (gridWidthNoPadding - (cellsPerRow * this.props.cellWidth)) / cellsPerRow);
        return {
            width: fittedWidth,
            height: fittedWidth
        };
    },

    componentDidMount: function() {
        var dimensions = this.calcCellDimensions();
        this.setState({
            cellWidth: dimensions.width,
            cellHeight: dimensions.height
        });
    },
    
    cellClicked: function(evt) {
        var $this = $(ReactDOM.findDOMNode(this));

        //width() returns without padding
        var gridWidthNoPadding = $this.width();
        
        var cellsPerRow = Math.floor(gridWidthNoPadding / this.state.cellWidth);

        var $target = $(evt.target);
        var targetPos = $target.position();

        var cellXPos = Math.floor((targetPos.left) / this.state.cellWidth);
        var yOffset = targetPos.top;

        var cellIndex = $target.data('cell-index')
        if (this.state.expanded) {

            //console.log('expanded cellindex: ' + cellIndex);
            if (cellIndex > (this.state.expanderIndex - 1)) {
                var expanderHeight = $this.find('.expander').height();
                yOffset -= expanderHeight;
            }
        }

        var cellYPos = Math.floor(yOffset / this.props.cellHeight);
        var expanderIndex = (cellYPos + 1) * cellsPerRow;

        //console.log('x:' + cellXPos + ', y:' + cellYPos);
        //console.log('expanderIndex: ' + expanderIndex);
        //console.log('cellsPerRow: ' + cellsPerRow);

        if (cellXPos === this.state.cellIndices.x &&
            cellYPos === this.state.cellIndices.y) {
            this.setState({
                expanded: false,
            });
        }
        else {
            this.setState({
                cellIndices: { x: cellXPos, y: cellYPos },
                expanderIndex: expanderIndex,
                selectedIndex: cellYPos * cellsPerRow + cellXPos,
                expanded: true,
                expanderContent: this.props.cells[cellIndex].content
            });
        }
    },

    expanderUnmounted: function() {
        // Once the expander has been unmounted, we can now indicate that
        // we want to group all the cells back together without the
        // expander split
        this.setState({
            cellIndices: { x:-1, y:-1},
            expanderIndex: -1,
            selectedIndex: -1,
        });
    },
    
    render: function() {
        function makeCellWrapper(index, selectedIndex, cell) {
            return (
                <div
                    key={index}
                    ref={"cellWrapper-" + index}
                    onClick={this.cellClicked}
                    className="cellWrapper pull-left"
                    data-cell-index={index}
                    style={{
                        width: this.state.cellWidth,
                        height: this.state.cellHeight,
                        //marginRight: this.props.spacingH,
                        //marginBottom: this.props.spacingV
                    }}>
                    {cell.cell}
                    <i className={ClassNames({
                                 "fa": true,
                                 "fa-caret-up": true,
                                 "hidden": index !== selectedIndex})}></i>
                </div>
            );
        }

        var beforeCells = [];
        var afterCells = [];
        for (var i=0; i<this.props.cells.length; ++i) {
            if (this.state.expanderIndex === -1 || i<this.state.expanderIndex) {
                beforeCells.push(makeCellWrapper.bind(this)(i, this.state.selectedIndex, this.props.cells[i]));
            } else {
                afterCells.push(makeCellWrapper.bind(this)(i, this.state.selectedIndex, this.props.cells[i]));
            }
        }

        var expander;
        var items = [];
        if (this.state.expanded) {
            items.push(<ExpanderWrapper
                           key="1"
                           unmounted={this.expanderUnmounted}
                           content={this.state.expanderContent}/>);
        }
                                                     
        expander = (
                    <ReactCSSTransitionGroup
                        transitionName="expander-animation"
                        transitionEnterTimeout={250}
                        transitionLeaveTimeout={250}>
                        {items}
                    </ReactCSSTransitionGroup>
        );

        return (
            <div className="cmp-Grid" style={{
                paddingLeft: this.props.paddingLeft,
                paddingRight: this.props.paddingRight,
                paddingTop: this.props.paddingTop,
                paddingBottom: this.props.paddingBottom,
            }}>
                <div className="beforeExpander">
                    {beforeCells}
                    <div style={{clear:"both"}}></div>
                </div>
                {expander}
                <div className="afterExpander">
                    {afterCells}
                    <div style={{clear:"both"}}></div>
                </div>
            </div>
        );
    }
});
module.exports = Grid;
