"use client";

import React, { useState, ChangeEvent } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

// 型定義
interface ChunkConfig {
  maxTokens: number;
  overlapTokens: number;
  minTokens: number;
  paragraphSeparator: string;
}

interface Chunk {
  content: string;
  startChar: number;
  endChar: number;
  tokenCount: number;
  index: number;
}

// サンプルテキスト（省略）
const SAMPLE_TEXT = `# 東京国際工科専門職大学について

東京国際工科専門職大学（IPUT）は、最先端のテクノロジーと実践的な専門教育を提供する専門職大学です。

## 教育理念

私たちは、技術革新と創造性を重視し、グローバルな視点を持つIT人材の育成に力を入れています。
産業界との密接な連携により、実践的なスキルと理論的知識の両方を身につけることができます。

## カリキュラム特徴

- プロジェクトベースの学習
- 第一線で活躍する実務家教員による指導
- 充実した英語教育プログラム

### 実践的な学び

1年次から実践的なプロジェクトに参加し、実際の課題解決に取り組みます。
企業との共同プロジェクトも多数実施しています。`;

const ChunkVisualizer: React.FC = () => {
  const [inputText, setInputText] = useState<string>(SAMPLE_TEXT);
  const [chunks, setChunks] = useState<Chunk[]>([]);
  const [config, setConfig] = useState<ChunkConfig>({
    maxTokens: 256,
    overlapTokens: 30,
    minTokens: 50,
    paragraphSeparator: "\n\n",
  });

  // 数値入力のハンドラー
  const handleNumberChange = (e: ChangeEvent<HTMLInputElement>, field: keyof ChunkConfig) => {
    const value = e.target.value;
    // 空文字列の場合は0をセット
    const numValue = value === "" ? 0 : Math.max(0, parseInt(value) || 0);
    setConfig((prev) => ({
      ...prev,
      [field]: numValue,
    }));
  };

  // 簡易的なトークンカウント（実際のGo実装を模倣）
  const countTokens = (text: string): number => {
    return text.split(/[\s\p{P}]+/u).filter(Boolean).length;
  };

  // チャンク分割処理（Go実装の疑似実装）
  const chunkDocument = (content: string): Chunk[] => {
    content = content.trim();
    if (!content) return [];

    const chunks: Chunk[] = [];
    const paragraphs = content.split(config.paragraphSeparator);
    let currentPos = 0;
    let nextStartPos = 0;

    paragraphs.forEach((para) => {
      para = para.trim();
      if (!para) {
        nextStartPos += config.paragraphSeparator.length;
        return;
      }

      currentPos = nextStartPos;
      const chunk: Chunk = {
        content: para,
        startChar: currentPos,
        endChar: currentPos + para.length,
        tokenCount: countTokens(para),
        index: chunks.length,
      };

      chunks.push(chunk);
      nextStartPos = chunk.endChar + config.paragraphSeparator.length;
    });

    return chunks;
  };

  const handleProcess = (): void => {
    const processedChunks = chunkDocument(inputText);
    setChunks(processedChunks);
  };

  return (
    <div className="space-y-8 p-6">
      <Card className="p-6">
        <div className="space-y-4">
          <h2 className="text-2xl font-bold">Chunk Visualizer</h2>

          {/* 設定部分 */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="maxTokens">Max Tokens</Label>
              <Input
                id="maxTokens"
                type="number"
                min={0}
                value={config.maxTokens}
                onChange={(e) => handleNumberChange(e, "maxTokens")}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="overlapTokens">Overlap Tokens</Label>
              <Input
                id="overlapTokens"
                type="number"
                min={0}
                value={config.overlapTokens}
                onChange={(e) => handleNumberChange(e, "overlapTokens")}
              />
            </div>
          </div>

          {/* 入力エリア */}
          <div className="space-y-2">
            <Label htmlFor="input">Input Text</Label>
            <Textarea
              id="input"
              className="h-64"
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
            />
          </div>

          <Button onClick={handleProcess}>Process Chunks</Button>
        </div>
      </Card>

      {/* チャンク可視化エリア */}
      <Card className="p-6">
        <h3 className="text-xl font-bold mb-4">Chunks Preview</h3>
        <div className="space-y-4">
          {chunks.map((chunk, index) => (
            <div key={index} className="border rounded-lg p-4 space-y-2">
              <div className="flex justify-between text-sm text-gray-500">
                <span>Chunk {index + 1}</span>
                <span>{chunk.tokenCount} tokens</span>
              </div>
              <div className="whitespace-pre-wrap">{chunk.content}</div>
              <div className="text-sm text-gray-500">
                Position: {chunk.startChar} - {chunk.endChar}
              </div>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
};

export default ChunkVisualizer;
