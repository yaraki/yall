#!/usr/bin/env -S npx ts-node

interface YExpr {
  print(): string
  equal(other: YExpr): boolean
  eval(env: YEnv): YExpr
}

abstract class YLiteral implements YExpr {
  abstract print(): string
  abstract equal(other: YExpr): boolean
  eval(env: YEnv): YExpr {
    return this;
  }
}

class YNumber extends YLiteral {
  constructor(readonly value: number) { super(); }
  print(): string {
    return `${this.value}`;
  }
  equal(other: YExpr): boolean {
    return other instanceof YNumber && this.value === other.value;
  }
}

class YBoolean extends YLiteral {
  static readonly TRUE = new YBoolean('#t');
  static readonly FALSE = new YBoolean('#f');

  private constructor(private readonly name: string) { super(); }
  print(): string {
    return this.name;
  }
  equal(other: YExpr): boolean {
    return this === other;
  }
}

class YString extends YLiteral {
  constructor(private readonly value: string) { super(); }
  print(): string {
    return `"${this.value}"`;
  }
  equal(other: YExpr): boolean {
    return other instanceof YString && this.value === other.value;
  }
}

class YSymbol implements YExpr {
  constructor(readonly name: string) { }
  print(): string {
    return this.name;
  }
  equal(other: YExpr): boolean {
    return other instanceof YSymbol && this.name === other.name;
  }
  eval(env: YEnv): YExpr {
    return env.resolve(this.name);
  }
}

class YQuoted implements YExpr {
  constructor(private readonly expr: YExpr) { }
  print(): string {
    return `'${this.expr.print()}`
  }
  equal(other: YExpr): boolean {
    return other instanceof YQuoted && this.expr.equal(other.expr);
  }
  eval(env: YEnv): YExpr {
    return this.expr;
  }
}

interface YList extends YExpr {
  printWithoutParens(): string;
  forEach(body: (expr: YExpr) => void): void;
  evalEach(env: YEnv): YList;
  nth(index: number): YExpr;
}

class YEmpty extends YLiteral implements YList {
  static EMPTY = new YEmpty();
  private constructor() { super(); }
  print(): string { return '()'; }
  printWithoutParens(): string { return ''; }
  equal(other: YExpr): boolean { return other instanceof YEmpty; }
  forEach(body: (expr: YExpr) => void): void { }
  evalEach(env: YEnv): YList { return this; }
  nth(index: number): YExpr {
    throw new Error('Index out of bounds');
  }
}

class YCell implements YList {

  constructor(
    readonly car: YExpr,
    readonly cdr: YList,
  ) { }

  print(): string {
    return `(${this.printWithoutParens()})`;
  }

  printWithoutParens(): string {
    if (this.cdr === YEmpty.EMPTY) {
      return this.car.print();
    }
    return `${this.car.print()} ${this.cdr.printWithoutParens()}`
  }

  equal(other: YExpr): boolean {
    return other instanceof YCell && this.car.equal(other.car) && this.cdr.equal(other.cdr);
  }

  eval(env: YEnv): YExpr {
    const callable = this.car.eval(env) as YCallable;
    if (isFunction(callable.body)) {
      return callable.body(this.cdr.evalEach(env));
    } else {
      return callable.body(env, this.cdr);
    }
  }

  evalEach(env: YEnv): YList {
    return new YCell(this.car.eval(env), this.cdr.evalEach(env));
  }

  forEach(body: (expr: YExpr) => void): void {
    body(this.car);
    this.cdr.forEach(body);
  }

  nth(index: number): YExpr {
    if (index == 0) {
      return this.car;
    }
    return this.cdr.nth(index - 1);
  }
}

type YFunctionBody = (args: YList) => YExpr;
type YFormBody = (env: YEnv, args: YList) => YExpr;
type YCallableBody = YFunctionBody | YFormBody;

function isFunction(body: YCallableBody): body is YFunctionBody {
  return body.length === 1;
}

enum CallableType { Function, Form };

class YCallable extends YLiteral {
  constructor(
    readonly name: string,
    readonly body: YCallableBody,
  ) { super(); }
  get type(): CallableType {
    if (this.body.arguments.length == 1) {
      return CallableType.Function;
    } else {
      return CallableType.Form;
    }
  }
  print(): string {
    return `#<FUNCTION ${this.name}>`;
  }
  equal(other: YExpr): boolean {
    return other instanceof YCallable && this.name === other.name && this.body === other.body;
  }
}

class YEnv {
  private readonly table: Map<string, YExpr> = new Map();

  constructor(private readonly parent?: YEnv | undefined) {
    for (const callable of YEnv.builtins) {
      this.table.set(callable.name, callable);
    }
  }

  private static readonly builtins = [
    new YCallable('if', (env: YEnv, args: YList) => {
      if (args instanceof YCell) {
        if (args.car.eval(env) !== YBoolean.FALSE) {
          return (args.cdr as YCell).car.eval(env);
        }
        return ((args.cdr as YCell).cdr as YCell).car.eval(env);
      }
      throw new Error('Empty if');
    }),
    new YCallable('def', (env: YEnv, args: YList) => {
      const symbol = args.nth(0) as YSymbol;
      const value = args.nth(1).eval(env);
      env.table.set(symbol.name, value);
      return symbol;
    }),
    new YCallable('fn', (env: YEnv, args: YList) => {
      const cell = args as YCell;
      const params = cell.car as YList;
      const body = cell.cdr;
      return new YCallable('#fn', (args: YList) => {
        const derived = new YEnv(env);
        derived.bindParams(params, args);
        return derived.begin(body);
      });
    }),
    new YCallable('+', (args: YList) => {
      let result = 0;
      args.forEach((expr: YExpr) => {
        if (expr instanceof YNumber) {
          result += expr.value;
        }
      });
      return new YNumber(result);
    }),
    new YCallable('*', (args: YList) => {
      let result = 1;
      args.forEach((expr: YExpr) => {
        if (expr instanceof YNumber) {
          result *= expr.value;
        }
      });
      return new YNumber(result);
    }),
    new YCallable('-', (args: YList) => {
      if (args instanceof YCell) {
        if (args.cdr instanceof YEmpty) {
          return new YNumber((args.car as YNumber).value * -1);
        }
        let result = (args.car as YNumber).value;
        args.cdr.forEach(expr => {
          result -= (expr as YNumber).value;
        })
        return new YNumber(result);
      }
      return new YNumber(0);
    }),
    new YCallable('/', (args: YList) => {
      if (args instanceof YCell) {
        if (args.cdr instanceof YEmpty) {
          return new YNumber(1 / (args.car as YNumber).value);
        }
        let result = (args.car as YNumber).value;
        args.cdr.forEach(expr => {
          result /= (expr as YNumber).value;
        });
        return new YNumber(result);
      }
      throw new Error('Illegal number of arguments');
    }),
    new YCallable('list', (args: YList) => args),
    new YCallable('car', (args: YList) => {
      return ((args as YCell).car as YCell).car;
    }),
    new YCallable('cdr', (args: YList) => {
      return ((args as YCell).car as YCell).cdr;
    }),
  ];

  eval(input: string): YExpr {
    const exprs = Array.from(read(input));
    if (exprs.length == 0) {
      return YEmpty.EMPTY;
    }
    return exprs[exprs.length - 1].eval(this);
  }

  resolve(symbol: string): YExpr {
    const expr = this.table.get(symbol);
    if (expr) {
      return expr;
    }
    if (this.parent === undefined) {
      throw Error(`Unbound symbol: ${symbol}`);
    }
    return this.parent.resolve(symbol);
  }

  bindParams(params: YList, args: YList): void {
    let cell = args as YCell;
    params.forEach(expr => {
      this.table.set((expr as YSymbol).name, cell.car);
      cell = cell.cdr as YCell;
    });
  }

  begin(body: YList): YExpr {
    let result: YExpr = YEmpty.EMPTY;
    body.forEach(expr => {
      result = expr.eval(this);
    });
    return result;
  }
}

function* tokenize(input: string): Generator<string> {
  function nextString(input: string, start: number): string {
    let buffer = '"';
    let escaped = false;
    for (let i = start + 1; i < input.length; ++i) {
      const c = input[i];
      if (escaped) {
        escaped = false;
      } else if (c == '\\') {
        escaped = true;
      } else if (c == '"') {
        buffer += '"';
        return buffer;
      }
      buffer += c;
    }
    throw new Error(`Malformed string ${buffer}`);
  }
  let buffer = '';
  for (let i = 0; i < input.length; ++i) {
    const c = input[i];
    switch (c) {
      case ' ':
      case '\n':
      case '\r':
      case '\t':
        if (buffer.length > 0) {
          yield buffer;
          buffer = '';
        }
        break;
      case '(':
      case ')':
      case '[':
      case ']':
      case "'":
      case '`':
        if (buffer.length > 0) {
          yield buffer;
          buffer = '';
        }
        yield c;
        break;
      case '"':
        const s = nextString(input, i);
        i += s.length - 1;
        yield s;
        break;
      default:
        buffer += c;
        break;
    }
  }
  if (buffer.length > 0) {
    yield buffer;
  }
}

function createList(items: YExpr[]): YList {
  let result = YEmpty.EMPTY;
  for (const item of items.reverse()) {
    result = new YCell(item, result);
  }
  return result;
}

function* read(input: string): Generator<YExpr> {
  yield* readTokens(tokenize(input));
}

function* readTokens(tokens: Generator<string>): Generator<YExpr> {
  for (let iter = tokens.next(); !iter.done; iter = tokens.next()) {
    const token = iter.value;
    if (token === '#t') {
      yield YBoolean.TRUE;
    } else if (token === '#f') {
      yield YBoolean.FALSE;
    } else if (token.startsWith('"') && token.endsWith('"')) {
      yield new YString(token.substring(1, token.length - 1));
    } else if (token === '(') {
      yield createList(Array.from(readTokens(tokens)));
    } else if (token === ')') {
      return;
    } else if (token === "'") {
      const expr = readTokens(tokens).next();
      yield new YQuoted(expr.value);
    } else {
      const f = parseFloat(token);
      const i = parseInt(token);
      if (!isNaN(f)) {
        yield new YNumber(f);
      } else if (!isNaN(i)) {
        yield new YNumber(i);
      } else {
        yield new YSymbol(token);
      }
    }
  }
}

function testPrint() {
  const testCases = new Map<YExpr, string>([
    [YEmpty.EMPTY, '()'],
    [new YCell(new YNumber(3), YEmpty.EMPTY), '(3)'],
    [new YCell(new YNumber(5), new YCell(new YNumber(3), YEmpty.EMPTY)), '(5 3)'],
    [YBoolean.TRUE, '#t'],
    [YBoolean.FALSE, '#f'],
    [new YString('abc'), '"abc"'],
  ]);
  for (const [expr, expected] of testCases) {
    const actual = expr.print();
    if (actual !== expected) {
      console.log(`Expected ${expected}, but received ${actual}`);
    }
  }
}

function testTokenize() {
  const testCases = new Map([
    ['abc', ['abc']],
    ['a b c', ['a', 'b', 'c']],
    [' a   b       c        ', ['a', 'b', 'c']],
    ['()', ['(', ')']],
    ['(a b)', ['(', 'a', 'b', ')']],
    ["'(1)", ["'", '(', '1', ')']],
    ['"abc"', ['"abc"']],
    ['(a "bc")', ['(', 'a', '"bc"', ')']],
    ['""', ['""']],
    ['"\\""', ['"\\""']],
    ['"a\\"b"', ['"a\\"b"']],
  ]);
  for (let [input, expected] of testCases) {
    let failed = false;
    let i = 0;
    const actual: string[] = [];
    for (let token of tokenize(input)) {
      if (token != expected[i]) {
        failed = true;
      }
      actual.push(token);
      ++i;
    }
    if (i != expected.length) {
      failed = true;
    }
    if (failed) {
      console.log(`Expected [${expected}] for ${input}, but received ${actual}`);
    }
  }
}

function testRead() {
  const testCases = new Map<string, YExpr>([
    ['abc', new YSymbol('abc')],
    ['1', new YNumber(1)],
    ['1.5', new YNumber(1.5)],
    ['"abc"', new YString('abc')],
    ['#t', YBoolean.TRUE],
    ['#f', YBoolean.FALSE],
    ['(a)', new YCell(new YSymbol('a'), YEmpty.EMPTY)],
    ['(a 3)', new YCell(new YSymbol('a'), new YCell(new YNumber(3), YEmpty.EMPTY))],
    ['((a 3) 5)', createList([createList([new YSymbol('a'), new YNumber(3)]), new YNumber(5)])],
    ["'a", new YQuoted(new YSymbol('a'))],
  ]);
  for (let [input, expected] of testCases) {
    const result = Array.from(read(input));
    if (result.length != 1) {
      console.log(`Multiple results from ${input}, received ${result.length} values`);
    } else {
      const actual = result[0];
      if (!actual.equal(expected)) {
        console.log(`Expected ${expected.print()}, but received ${actual.print()}`);
      }
    }
  }
}

function testEval() {
  const testCases = new Map<string, string>([
    ['1', '1'],
    ['"a"', '"a"'],
    ["'a", 'a'],
    ['#t', '#t'],
    ['#f', '#f'],
    ['(+ 1 2)', '3'],
    ['(* 1 2 3 4 5)', '120'],
    ['(- 3)', '-3'],
    ['(- 10 2 3)', '5'],
    ['(/ 2)', '0.5'],
    ['(/ 1001 13 11)', '7'],
    ['(if #t 1 a)', '1'],
    ['(if 0 1 a)', '1'],
    ['(if #f a 1)', '1'],
    ["'(a b c)", '(a b c)'],
    ['(list (+ 1 2))', '(3)'],
    ["'(+ 1 2)", '(+ 1 2)'],
    ["(def a 3) '(a 5)", '(a 5)'],
    ["(def a 3) (list a 5)", '(3 5)'],
    ['((fn (x) (* x 2)) 3)', '6'],
    ['(def a (fn (x) (* x 2))) (a 3)', '6'],
    ['(car (list 1 2))', '1'],
    ['(cdr (list 1 2))', '(2)'],
  ]);
  const env = new YEnv();
  for (const [input, expected] of testCases) {
    let actual: YExpr = YEmpty.EMPTY;
    for (const expr of read(input)) {
      actual = expr.eval(env);
    }
    const actualPrint = actual.print();
    if (actualPrint !== expected) {
      console.log(`Expected ${expected}, but received ${actualPrint}`);
    }
  }
}

async function main() {
  const env = new YEnv();
  process.stdout.write('yall> ');
  for await (const chunk of process.stdin) {
    const input = chunk.toString();
    const expr = env.eval(input);
    console.log(expr.print());
    process.stdout.write('yall> ');
  }
}

if (process.argv[2] === '-t') {
  testPrint();
  testTokenize();
  testRead();
  testEval();
} else {
  main();
}
