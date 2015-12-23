(function() {

    /*
    var App = React.createClass({
        render: function() {
            return (
                <div>Welcome to gohome</div>
            );
        }
    });

    debugger;
    React.render(<App />, document.body);
     */
    return;
    setTimeout(function() {
        var el = document.getElementsByClassName("sceneList")[0];
        var sortable = Sortable.create(el);
    }, 2000);
})();
