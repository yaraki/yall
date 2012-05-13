#import('dart:html');

class ReaderException implements Exception {
  final String message;
  ReaderException(this.message);
}

class RuntimeException implements Exception {
  final String message;
  RuntimeException(this.message);
}

interface YExpr {
  String str();
  YExpr eval(YEnv env);
}

abstract class YLiteral implements YExpr {

  const YLiteral();

  YExpr eval(YEnv env) {
    return this;
  }
}

class YBoolean extends YLiteral {

  static final YBoolean TRUE = const YBoolean('#t');
  static final YBoolean FALSE = const YBoolean('#f');

  final String name;

  const YBoolean(this.name);

  String str() {
    return this.name;
  }
}

class YNumber extends YLiteral {

  final num value;

  YNumber(this.value) {
  }

  String str() {
    return "${this.value}";
  }
}

class YString extends YLiteral {

  final String value;

  YString(this.value) {
  }

  String str() {
    return "\"${this.value}\"";
  }
}

class YSymbol implements YExpr {

  final String value;

  YSymbol(this.value) {
  }

  String str() {
    return this.value;
  }

  YExpr eval(YEnv env) {
    return env.resolve(this);
  }
}

interface YList extends YExpr {
  String strWithoutParens();
}

class YEmpty extends YLiteral implements YList {

  static final YEmpty EMPTY = const YEmpty();

  const YEmpty();

  String strWithoutParens() {
    return "";
  }

  String str() {
    return "()";
  }
}

class YCell implements YList {

  YExpr head;
  YList tail;

  YCell(YExpr head, YCell tail) {
    this.head = head;
    this.tail = tail;
  }

  String strWithoutParens() {
    if (this.tail is YEmpty) {
      return this.head.str();
    }
    return "${this.head.str()} ${this.tail.strWithoutParens()}";
  }

  String str() {
    return "(${strWithoutParens()})";
  }

  YExpr eval(YEnv env) {
    YExpr expr = env.eval(this.head);
    switch (true) {
      case expr is YFunction:
        YFunction func = expr;
        return func.body(env.evalEach(this.tail));
    }
    throw new RuntimeException("Cannot eval: ${this.str()}");
  }
  
  void forEach(void func(YExpr expr)) {
    func(this.head);
    if (this.tail is YCell) {
      YCell cell = this.tail;
      cell.forEach(func);
    }
  }
}

typedef YExpr YFunctionBody(YList args);

class YFunction extends YLiteral {

  final String name;
  final YFunctionBody body;
  
  YFunction(String this.name, YFunctionBody this.body);

  String str() {
    return "#<FUNCTION ${this.name}>";
  }
}

class YReader {

  String input;
  int position;

  YReader(String input) {
    this.input = input;
    this.position = 0;
  }

  String nextString() {
    ++this.position;  // Skip the first '"'
    StringBuffer sb = new StringBuffer('"');
    bool escaped = false;
    int max = this.input.length;
    for (; this.position < max; ++this.position) {
      var c = this.input[this.position];
      if (escaped) {
        escaped = false;
      } else if (c == '\\') {
        escaped = true;
      } else if (c == '"') {
        ++this.position;
        sb.add('"');
        return sb.toString();
      }
      sb.add(c);
    }
    throw new ReaderException("Illegal end of String.");
  }

  String nextToken() {
    StringBuffer sb = new StringBuffer();
    int max = this.input.length;
    for (; this.position < max; ++this.position) {
      String char = this.input[this.position];
      switch (char) {
        case ' ':
        case '\n':
        case '\r':
        case '\t':
          if (0 < sb.length) {
            ++this.position;
            return sb.toString();
          }
          break;
        case '(':
        case ')':
        case '[':
        case ']':
        case '\'':
        case '`':
          if (0 < sb.length) {
            return sb.toString();
          }
          ++this.position;
          return char;
        case '"':
          return nextString();
        default:
          sb.add(char);
          break;
      }
    }
    if (0 < sb.length) {
      return sb.toString();
    }
    return null;
  }

  YExpr readExpr(bool asList) {
    String token = nextToken();
    print("token: ${token}");
    YExpr retval;
    switch (token) {
      case '(': // list start
        retval = readExpr(true);
        break;
      case ')': // list end
        retval = YEmpty.EMPTY;
        asList = false;
        break;
      default:
        if (token[0] == '"' && token[token.length - 1] == '"') { // YString
          retval = new YString(token.substring(1, token.length - 1));
        } else { // YNumber
          try {
            retval = new YNumber(Math.parseInt(token));
          } catch (BadNumberFormatException e1) {
            try {
              retval = new YNumber(Math.parseDouble(token));
            } catch (BadNumberFormatException e2) {
              // YSymbol
              retval = new YSymbol(token);
            }
          }
        }
        break;
    }
    if (asList) {
      return new YCell(retval, readExpr(true));
    }
    return retval;
  }

  YExpr read() {
    return readExpr(false);
  }
}

class YEnv {

  Map table;

  YEnv() {
    this.table = new Map();
    bindFunction('+', (YList args) {
      YCell cell = args;
      num result = 0;
      cell.forEach((YExpr expr) {
        YNumber number = expr;
        result += number.value;
      });
      return new YNumber(result);
    });
    bindFunction('-', (YList args) {
      YCell cell = args;
      num result = 0;
      bool first = true;
      cell.forEach((YExpr expr) {
        YNumber number = expr;
        if (first) {
          first = false;
          result += number.value;
        } else {
          result -= number.value;
        }
      });
      return new YNumber(result);
    });
    bindFunction('*', (YList args) {
      YCell cell = args;
      num result = 1;
      cell.forEach((YExpr expr) {
        YNumber number = expr;
        result *= number.value;
      });
      return new YNumber(result);
    });
    bindFunction('/', (YList args) {
      YCell cell = args;
      num result = 1;
      bool first = true;
      cell.forEach((YExpr expr) {
        YNumber number = expr;
        if (first) {
          first = false;
          result *= number.value;
        } else {
          result /= number.value;
        }
      });
      return new YNumber(result);
    });
  }
  
  void bindFunction(String name, YFunctionBody body) {
    bind(new YSymbol(name), new YFunction(name, body));
  }

  void bind(YSymbol symbol, YExpr expr) {
    table[symbol.value] = expr;
  }

  YExpr resolve(YSymbol symbol) {
    if (!table.containsKey(symbol.value)) {
      throw new RuntimeException("Unbound symbol: ${symbol.value}");
    }
    return table[symbol.value];
  }
  
  YList evalEach(YList list) {
    if (list is YEmpty) {
      return YEmpty.EMPTY;
    }
    YCell cell = list;
    return new YCell(eval(cell.head), evalEach(cell.tail));
  }

  YExpr eval(YExpr expr) {
    return expr.eval(this);
  }
}

void main() {
  document.query('#status').innerHTML = 'Ready.';
  YEnv env = new YEnv();
  document.query('#run').on.click.add((event) {
    InputElement input = document.query('#input');
    YReader reader = new YReader(input.value);
    Element readResult = document.query('#read_result');
    Element evalResult = document.query('#eval_result');
    readResult.classes.clear();
    evalResult.classes.clear();
    try {
      YExpr expr = reader.read();
      readResult.innerHTML = expr.str();
      evalResult.innerHTML = env.eval(expr).str();
    } catch (ReaderException e) {
      readResult.innerHTML = e.message;
      readResult.classes.add("Error");
    } catch (RuntimeException e) {
      evalResult.innerHTML = e.message;
      evalResult.classes.add("Error");
    }
  });
}
