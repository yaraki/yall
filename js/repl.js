(function() {
    function println(message) {
        $('#output').append($('<li>').text(message));
    }

    $(function() {
        var env = new yall.Env();
        println(env.evalString('12').print());
        println('It is fine today');
    });
})();