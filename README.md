# IPUT TOKYO AI

東京国際工科専門職大学（IPUT）のためのRAGシステムを備えたQ&A ChatBot

## 開発中...

現在、リリースに向けて、IPUTの学生である私が個人で開発を進めています🚀

GDG Devfes Tokyo 2024で登壇予定です！

## curlコマンド例

質問をする
```
curl -X POST http://localhost:9020/query/ -H "Content-Type: application/json" -d '{"content": "情報工学科について教えてください"}'
```

ドキュメントを追加する
```
curl -X POST http://localhost:9020/add/ -H "Content-Type: application/json" -d @server/university_data.json
```
