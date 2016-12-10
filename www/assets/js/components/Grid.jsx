var React = require('react');
var ReactTransitionGroup = require('react-addons-transition-group');
var ReactDOM = require('react-dom');
var ClassNames = require('classnames');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'Grid',
    prefix: 'b-'
});
require('../../css/components/Grid.less')

var Grid = React.createClass({
    debug: function(msg) {
        if (this.props.debugName) {
            console.log('[' + this.props.debugName + '] - ' + msg);
        }
    },

    getDefaultProps: function() {
        return {
            cells: [],
            cellWidth: 110,
            cellHeight: 125,
            spacingH: 0,
            spacingV: 0
        };
    },

    getInitialState: function() {
        this.debug('getInitialState');

        return {
            cellIndices: { x:-1, y:-1 },
            expanderIndex: -1,
            selectedIndex: -1,
            expanded: false,
            cellWidth: this.props.cellWidth,
            cellHeight: this.props.cellHeight
        }
    },

    componentWillReceiveProps: function(nextProps) {
        // If the cells changed, then we are goign to re-render, close the expander
        if (nextProps.cells && (nextProps.cells !== this.props.cells)) {
            //TODO: Needed? causing issues when updating items on import, revist
            //this.closeExpander();
        }
    },

    shouldComponentUpdate: function(nextProps, nextState) {
        //TODO: Fix
        return true;
        /*
        if (nextProps.cells && (nextProps.cells != this.props.cells)) {
            return true;
        }
        if (nextState.cellWidth && (nextState.cellWidth !== this.state.cellWidth)) {
            return true;
        }
        if (nextState.cellHeight && (nextState.cellHeight !== this.state.cellHeight)) {
            return true;
        }
        if (nextState.expanderIndex != undefined && (nextState.expanderIndex !== this.state.expanderIndex)) {
            return true;
        }
        return false;*/
    },

    calcCellDimensions: function() {
        var $this = $(ReactDOM.findDOMNode(this));
        var gridWidthNoPadding = $this.width();
        var cellsPerRow = Math.floor(gridWidthNoPadding / this.props.cellWidth);

        var fittedWidth = Math.floor(this.props.cellWidth + (gridWidthNoPadding - (cellsPerRow * this.props.cellWidth)) / cellsPerRow);
        return {
            width: fittedWidth,
            // height always remains constant
            height: this.props.cellHeight
        };
    },

    componentDidMount: function() {
        if (this.props.debugName) {
            console.log('mounting grid: ' + this.props.debugName);
        }

        var dimensions = this.calcCellDimensions();
        this.setState({
            cellWidth: dimensions.width,
            cellHeight: dimensions.height
        });
    },

    componentDidUpdate: function() {
        var dimensions = this.calcCellDimensions();
        if (dimensions.width !== this.state.cellWidth ||
            dimensions.height !== this.state.cellHeight) {
            this.setState({
                cellWidth: dimensions.width,
                cellHeight: dimensions.height
            });
        }
    },

    closeExpander: function() {
        this.debug('closeExpander');

        this.setState({
            expanded: false,
            cellIndices: { x:-1, y:-1},
            expanderIndex: -1,
            selectedIndex: -1,
            expanderContent: null
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
            // Have to take into account the expander height when calculating which
            // cell the user is clicking on
            if (cellIndex > (this.state.expanderIndex - 1)) {
                var expanderHeight = $this.find('.b-Expander').height();
                yOffset -= expanderHeight;
            }
        }

        var cellYPos = Math.floor(yOffset / this.props.cellHeight);
        var expanderIndex = Math.min(this.props.cells.length, (cellYPos + 1) * cellsPerRow);

        if (cellXPos === this.state.cellIndices.x &&
            cellYPos === this.state.cellIndices.y) {
            this.closeExpander();
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

    expanderWillMount: function(content) {
        this.props.expanderWillMount && this.props.expanderWillMount(content);
    },

    render: function() {
        function makeCellWrapper(index, selectedIndex, cell) {
            var content = cell.cell;
            return (
                <div
                    key={cell.key}
                    ref={"cellWrapper-" + index}
                    onClick={this.cellClicked}
                    {...classes('cell', '', 'pull-left')}
                    data-cell-index={index}
                    style={{
                        width: this.state.cellWidth,
                        height: this.state.cellHeight,
                    }}>
                    {content}
                    <i {...classes('expanded-arrow', index !== selectedIndex ? 'hidden' : '', 'fa fa-caret-up')}></i>
                </div>
            );
        }

        var content = [];
        var key = '';
        if (this.state.selectedIndex !== -1) {
            key = this.props.cells[this.state.selectedIndex].key;
        }

        if (this.props.debugName) {
            console.log('[' + this.props.debugName + '] selectedIndex: ' + this.state.selectedIndex);
        }

        var transitionGroup = (
            <ReactTransitionGroup key={key + "transition"}>
                <ExpanderWrapper key={key + 'wrapper'}>
                    {this.state.expanderContent}
                </ExpanderWrapper>
            </ReactTransitionGroup>
        );

        for (var i=0; i<this.props.cells.length; ++i) {
            content.push(makeCellWrapper.bind(this)(i, this.state.selectedIndex, this.props.cells[i]));
            if (this.state.expanded && this.state.expanderIndex === (i + 1)) {
                content.push(transitionGroup);
            }
        }

        if (!this.state.expanded) {
            //TODO: This isn't working, looks like multiple transition groups are causing issues, revist.
            //This should make the expander animate closed, but seems to have some bugs
            //content.push(transitionGroup);
        }

        return (
            <div {...classes()}>
                <div className="clearfix beforeExpander">
                    {content}
                    <div style={{clear:"both"}}></div>
                </div>
            </div>
        );
    }
});

var expClasses = new BEMHelper({
    name: 'Expander',
    prefix: 'b-'
});

var ExpanderWrapper = React.createClass({
    // Part of the calls for TransitionGroup, have to set the CSS property
    // to the initial value, then after a small delay set the end animation value
    componentWillAppear: function(cb) {
        var $this = $(ReactDOM.findDOMNode(this)).find('.animateWrapper');
        $this.css({ 'margin-top': -500 });
        setTimeout(function() {
            cb();
        }, 10);
    },
    componentDidAppear: function() {
        var $this = $(ReactDOM.findDOMNode(this)).find('.animateWrapper');
        $this.css({ 'margin-top': '0px' });
    },
    componentWillLeave: function(cb) {
        var $this = $(ReactDOM.findDOMNode(this)).find('animateWrapper');
        $this.css({ 'margin-top': -500 });
        cb();
    },
    render: function() {
        return (
            <div {...expClasses()}>
                <div className="animateWrapper">
                    {this.props.children}
                </div>
            </div>
        );
    }
});

module.exports = Grid;
