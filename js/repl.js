(function() {

    var env = new yall.Env();

    var log = {};
    log.normal = function(message) {
        $('#output').prepend($('<li>').text(message).addClass('normal'));
    };
    log.error = function(message) {
        $('#output').prepend($('<li>').text(message).addClass('error'));
    };

    $(function() {
        var $input = $('#input');
        $input.keypress(function(e) {
            if ((e.which && e.which == 13) || (e.keyCode && e.keyCode == 13)) { // enter
                var text = $input.val();
                $input.val('');
                if (/^\s*$/.test(text)) {
                    return false;
                }
                try {
                    log.normal(env.evalString(text).print());
                } catch (x) {
                    log.error(x);
                }
                return false;
            }
            return true;
        });
    });
})();