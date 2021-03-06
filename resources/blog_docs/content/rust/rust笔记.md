## 安装

注意！！！
-	Rust 风格的缩进使用 4 个空格，而不是制表符;
-	函数后加!代表调用 Rust 宏 ，如println! 而不是普通的函数;
-	Rust 代码的大多数行都以一个 ; 结尾；



## 使用

#### rustc
-	根据rustc编译 main.rs 生成可执行文件 

#### Cargo 是 Rust 的构建系统和包管理器。
    cargo new name 新建项目
		如果在现有的 Git 仓库中运行 cargo new，则不会生成 Git 文件；可以使用 cargo new --vcs=git 来覆盖。
	cargo build 构建项目
	cargo run   构建并运行项目
	cargo check 校验项目，不生成可执行文件

## 发布
当项目最终准备好发布时，可以使用 cargo build --release 来优化编译项目。这会在 target/release 而不是 target/debug 下生成可执行文件。这些优化可以让 Rust 代码运行的更快，不过启用这些优化也需要消耗更长的编译时间。这也就是为什么会有两种不同的配置：一种是为了开发，你需要经常快速重新构建；另一种是为用户构建最终程序，它们不会经常重新构建，并且希望程序运行得越快越好。如果你要对代码运行时间进行基准测试，请确保运行 cargo build --release 并使用 target/release 下的可执行文件进行测试。




## 变量
	let apples = 5; // 不可变
	let mut bananas = 5; // 可变
	let mut guess = String::new();创建了一个可变变量，并绑定到一个新的String空实例上;

const(注明类型) 常量
let 声明变量(mut 可变)

#### 遮蔽
	let x = 5;
	let x = x + 1;
	二次声明变量会覆盖变量值和类型，而mut只改变值，不改变类型;


## 数据类型

#### 标量（scalar）表示单个值
Rust 有 4 个基本的标量类型：整型、浮点型、布尔型和字符。
    整数（integer）是没有小数部分的数字。我们在第 2 章使用过一个整数类型（整型），即 u32 类型。此类型声明表明它关联的值应该是占用 32 位空间的无符号整型（有符号整型以 i 开始，i 是英文单词 integer 的首字母，与之相反的是 u，代表无符号 unsigned 类型）(i/u 8 16 32 64 128 size)，isize 和 usize 类型取决于程序运行的计算机体系结构，在表中表示为“arch”：若使用 64 位架构系统则为 64 位，若使用 32 位架构系统则为 32 位。

| 长度   | 有符号类型 | 无符号类型 |
| ------ | ---------- | ---------- |
| 8 位   | `i8`       | `u8`       |
| 16 位  | `i16`      | `u16`      |
| 32 位  | `i32`      | `u32`      |
| 64 位  | `i64`      | `u64`      |
| 128 位 | `i128`     | `u128`     |
| arch   | `isize`    | `usize`    |

整型溢出:
- 1、当在调试（debug）模式编译时，Rust 会检查整型溢出，若存在这些问题则使程序在编译时 panic； 
- 2、当使用 --release 参数进行发布（release）模式构建时，Rust 不检测会导致 panic 的整型溢出，大于该类型最大值的数值会被“包裹”成该类型能够支持的对应数字的最小值，如u8的257 会显示成1。