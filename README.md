# ccbuddy(ClaudeCode Buddy)

cceggはClaudeCodeを使うごとに成長するTUIで表現されるキャラクターです
現在のREADME.mdはこれから作成されるプロジェクトの要件などを示したものです

## 用語

- ccbuddyで表されるキャラクターはbuddyと呼ばれます

## 要件

- buddyには見た目、年齢、健康状態というパラメータがあります
- buddyは最初は卵の状態ですが、ユーザーがClaudeCodeを少し使うことで孵化しランダムなキャラクターが生成されます
- 孵化以降、ユーザーがClaudeCodeを使うと自動的にBuddyに対して餌が与えられます
- ユーザーがClaudeCodeを使っていないとBuddyには餌が与えられないため健康状態が悪くなり最悪の場合死にます
- よって、リアルタイムで成長するため定期的なケアが必要となります
- ユーザーはccbuddyとコマンドを叩くことでbuddyを呼び出すができます

## 技術要素

bubbleteaを使用します
