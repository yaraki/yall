(function () {

    var yall = {
        Cell: function(car, cdr) {
            this.init(car, cdr);
        },
        EMPTY: {
            print: function() { return '()'; },
            applyEval: function(env) { return this; }
        },
        Symbol: function(name) {
            this.init(name);
        },
        String: function(value) {
            this.init(value);
        },
        Number: function(value) {
            this.init(value);
        },
        Function: function(name, type, func) {
            this.init(name, type, func);
        },
        TRUE: {
            print: function () { return '#t'; },
            applyEval: function(env) { return this; }
        },
        FALSE: {
            print: function () { return '#f'; },
            applyEval: function(env) { return this; }
        },
        FUNC_PROC: 1,
        FUNC_FORM: 2,
        FUNC_MACRO: 3,
        Reader: function(s) {
            this.init(s);
        },
        Env: function(origin) {
            this.init(origin);
        }
    };

    yall.Cell.prototype = {
        init: function(car, cdr) {
            this.car = car;
            this.cdr = cdr;
        },

        print: function() {
            return '(' + this.printWithoutParens() + ')';
        },

        printWithoutParens: function() {
            if (yall.EMPTY == this.cdr) {
                return this.car.print();
            }
            return this.car.print() + ' ' + this.cdr.printWithoutParens();
        },

        applyEval: function(env) {
            var func = env.eval(this.car);
            switch (func.type) {
              case yall.FUNC_PROC:
                return func.func(env.evalEach(this.cdr));
              case yall.FUNC_FORM:
                return func.func(env, this.cdr);
            }
            throw 'Invaid application: ' + this.print();
        }
    };

    yall.Symbol.prototype = {
        init: function(name) {
            this.name = name;
        },

        print: function() {
            return this.name;
        },

        applyEval: function(env) {
            return env.resolve(this.name);
        }
    };

    yall.String.prototype = {
        init: function(value) {
            this.value = value;
        },

        print: function() {
            return '"' + this.value + '"';
        },

        applyEval: function(env) {
            return this;
        }
    };

    yall.Number.prototype = {
        init: function(value) {
            this.value = value;
        },

        print: function() {
            return '' + this.value;
        },

        applyEval: function(env) {
            return this;
        }
    };

    var FUNCTION_TYPENAMES = [];
    FUNCTION_TYPENAMES[yall.FUNC_PROC] = 'proc';
    FUNCTION_TYPENAMES[yall.FUNC_FORM] = 'form';
    FUNCTION_TYPENAMES[yall.FUNC_MACRO] = 'macro';

    yall.Function.prototype = {
        init: function(name, type, func) {
            this.name = name;
            this.type = type;
            this.func = func;
        },

        print: function() {
            return '#<' + FUNCTION_TYPENAMES[this.type] + ' ' + this.name + '>';
        },

        applyEval: function(env) {
            return this;
        }
    };

    yall.Reader.prototype = {
        init: function(s) {
            this.input = s.split('');
            this.pos = 0;
        },

        nextString: function() {
            var buffer = '"';
            ++this.pos; // skip first '"' in this.input
            var escaped = false;
            for (var max = this.input.length; this.pos < max; ++this.pos) {
                var c = this.input[this.pos];
                if (escaped) {
                    escaped = false;
                } else if (c == '\\') {
                    escaped = true;
                } else if (c == '"') {
                    ++this.pos;
                    return buffer + '"';
                }
                buffer += c;
            }
            throw "Unexpected end of string";
        },

        nextToken: function() {
            var buffer = '';
            for (var max = this.input.length; this.pos < max; ++this.pos) {
                var c = this.input[this.pos];
                if (/\s/.test(c)) { // space
                    if (0 < buffer.length) {
                        ++this.pos;
                        return buffer;
                    }
                } else if (/[()'`\[\]]/.test(c)) {
                    if (0 < buffer.length) {
                        return buffer;
                    }
                    ++this.pos;
                    return c;
                } else if (',' == c) {
                    if (0 < buffer.length) {
                        return buffer;
                    }
                    if ('@' == this.input[this.pos + 1]) {
                        this.pos += 2;
                        return ',@';
                    }
                    ++this.pos;
                    return ',';
                } else if ('"' == c) {
                    return this.nextString();
                } else {
                    buffer += c;
                }
            }
            if (0 < buffer.length) {
                return buffer;
            }
            return null;
        },

        readToken: function(asList) {
            var token = this.nextToken();
            var ret = null;
            if (token == '(') {
                ret = this.readToken(true);
            } else if (token == ')') {
                ret = yall.EMPTY;
                asList = false;
            } else if (/^".*"$/.test(token)) {
                ret = new yall.String(token.slice(1, -1));
            } else {
                var i = parseInt(token), f = parseFloat(token);
                if (i) {
                    if (f) {
                        ret = new yall.Number(f);
                    } else {
                        ret = new yall.Number(i);
                    }
                } else {
                    ret = new yall.Symbol(token);
                }
            }
            if (asList) {
                return new yall.Cell(ret, this.readToken(true));
            }
            return ret;
        },

        read: function() {
            return this.readToken(false);
        }
    };

    var builtins = [
        // Built-in special forms

        new yall.Function('def', yall.FUNC_FORM, function (env, args) {
            var symbol = args.car;
            var value = env.eval(args.cdr.car);
            env.bind(symbol.name, value);
            return symbol;
        }),

        new yall.Function('if', yall.FUNC_FORM, function (env, args) {
            var condition = env.eval(args.car);
            if (condition != yall.FALSE) {
                return env.eval(args.cdr.car);
            }
            return env.eval(args.cdr.cdr.car);
        }),

        new yall.Function('eval', yall.FUNC_FORM, function (env, args) {
            args = env.evalEach(args);
            return env.eval(args.car);
        }),

        // Built-in functions

        new yall.Function('list', yall.FUNC_PROC, function (args) {
            return args;
        }),

        new yall.Function('car', yall.FUNC_PROC, function (args) {
            return args.car.car;
        }),

        new yall.Function('cdr', yall.FUNC_PROC, function (args) {
            return args.car.cdr;
        }),

        new yall.Function('null?', yall.FUNC_PROC, function (args) {
            if (args.car == yall.EMPTY) {
                return yall.TRUE;
            }
            return yall.FALSE;
        }),

        new yall.Function('+', yall.FUNC_PROC, function(args) {
            var result = 0;
            for (var cell = args; cell != yall.EMPTY; cell = cell.cdr) {
                result += cell.car.value;
            }
            return new yall.Number(result);
        }),

        new yall.Function('-', yall.FUNC_PROC, function (args) {
            if (args.cdr == yall.EMPTY) {
                return new yall.Number(args.car.value * -1);
            }
            var result = args.car.value;
            for (var cell = args.cdr; cell != yall.EMPTY; cell = cell.cdr) {
                result -= cell.car.value;
            }
            return new yall.Number(result);
        }),

        new yall.Function('*', yall.FUNC_PROC, function (args) {
            var result = 1;
            for (var cell = args; cell != yall.EMPTY; cell = cell.cdr) {
                result *= cell.car.value;
            }
            return new yall.Number(result);
        }),

        new yall.Function('/', yall.FUNC_PROC, function (args) {
            if (args.cdr == yall.EMPTY) {
                return new yall.Number(1 / args.car.value);
            }
            var result = args.car.value;
            for (var cell = args.cdr; cell != yall.EMPTY; cell = cell.cdr) {
                result /= cell.car.value;
            }
            return new yall.Number(result);
        })
    ];

    yall.Env.prototype = {

        init: function(origin) {
            this.origin = origin || null;
            this.symbols = {};
            for (var i = 0, l = builtins.length; i < l; ++i) {
                var f = builtins[i];
                this.bind(f.name, f);
            }
            this.bind(yall.TRUE.print(), yall.TRUE);
            this.bind(yall.FALSE.print(), yall.FALSE);
        },

        bind: function(symbolName, value) {
            this.symbols[symbolName] = value;
        },

        resolve: function(symbolName) {
            var value = this.symbols[symbolName];
            if (value) {
                return value;
            }
            if (this.origin) {
                return this.origin.resolve(symbolName);
            }
            throw 'Unbound symbol: ' + symbolName;
        },

        evalString: function(s) {
            var reader = new yall.Reader(s);
            return this.eval(reader.read());
        },

        eval: function(expr) {
            return expr.applyEval(this);
        },

        evalEach: function(list) {
            if (list == yall.EMPTY) {
                return yall.EMPTY;
            }
            return new yall.Cell(this.eval(list.car), this.evalEach(list.cdr));
        }
    };

    function testReaderNextToken() {
        var reader = new yall.Reader('(a (b) \'(c   d) `(e ,f ,@g) "h  \\"i"))');
        var tokens = [];
        var token;
        while (null != (token = reader.nextToken())) {
            tokens.push(token);
        }
        var expected = ['(', 'a', '(', 'b', ')', '\'', '(', 'c', 'd', ')',
                        '`', '(', 'e', ',', 'f', ',@', 'g', ')', '"h  \\"i"', ')', ')'];
        for (var i = 0, l = expected.length; i < l; ++i) {
            if (expected[i] != tokens[i]) {
                console.log('Received ' + tokens[i] + ", expected " + expected[i]);
            }
        }
    }

    function testReaderRead() {
        var s = '(a b ("c") 123 12.5)';
        var reader = new yall.Reader(s);
        var expr = reader.read();
        if (expr.print() != s) {
            console.log('Received ' + expr.print() + ', expected ' + s);
        }
    }

    var evalTestCases = [
        {
            settings: ['(def a 3)'],
            input: '(list 1 2 a)',
            expected: '(1 2 3)'
        },
        {
            settings: [],
            input: '(if #t 1 2)',
            expected: '1'
        },
        {
            settings: [],
            input: '(null? ())',
            expected: '#t'
        },
        {
            settings: [],
            input: '(null? 12)',
            expected: '#f'
        },
        {
            settings: ['(def a 5)'],
            input: '(eval a)',
            expected: '5'
        },
        {
            settings: [],
            input: '(+ 1 2)',
            expected: '3'
        }
/*
        {
            settings: [],
            input: '',
            expected: ''
        }
*/
    ];

    function testEnvEval() {
        for (var i = 0, l = evalTestCases.length; i < l; ++i) {
            var env = new yall.Env();
            var tc = evalTestCases[i];
            for (var j = 0, m = tc.settings.length; j < m; ++j) {
                env.evalString(tc.settings[j]);
            }
            var result = env.evalString(tc.input);
            if (result.print() != tc.expected) {
                console.log('Received ' + result.print() + ', expected ' + tc.expected);
            }
        }
    }

    testReaderNextToken();
    testReaderRead();
    testEnvEval();

    if ('undefined' != typeof(window)) {
        window.yall = yall;
    }
    return yall;
})();
