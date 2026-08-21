import { describe, expect, test } from "bun:test";

import { buildNodeMentionReferences, collectUpstreamVideoNodes } from "../src/lib/canvas/canvas-resource-references";
import { CanvasNodeType, type CanvasConnection, type CanvasNodeData } from "../src/types/canvas";

function videoNode(id: string): CanvasNodeData {
    return {
        id,
        type: CanvasNodeType.Video,
        title: id,
        position: { x: 0, y: 0 },
        width: 100,
        height: 100,
        metadata: { content: `data:video/mp4;base64,${id}` },
    };
}

function textNode(id: string): CanvasNodeData {
    return {
        id,
        type: CanvasNodeType.Text,
        title: id,
        position: { x: 0, y: 0 },
        width: 100,
        height: 60,
        metadata: { content: id },
    };
}

function connection(fromNodeId: string, toNodeId: string): CanvasConnection {
    return { id: `conn-${fromNodeId}-${toNodeId}`, fromNodeId, toNodeId };
}

function imageNode(id: string, metadata: CanvasNodeData["metadata"] = { content: `data:image/png;base64,${id}` }): CanvasNodeData {
    return {
        id,
        type: CanvasNodeType.Image,
        title: id,
        position: { x: 0, y: 0 },
        width: 100,
        height: 100,
        metadata,
    };
}

function configNode(id: string): CanvasNodeData {
    return {
        id,
        type: CanvasNodeType.Config,
        title: id,
        position: { x: 0, y: 0 },
        width: 100,
        height: 60,
        metadata: {},
    };
}

describe("buildNodeMentionReferences", () => {
    test("视频节点会同时展示连入的文本和图片预引用", () => {
        const video = videoNode("video");
        const text = textNode("note");
        const image = imageNode("still");
        const nodes = [video, text, image];
        const connections = [connection("note", "video"), connection("still", "video")];
        expect(buildNodeMentionReferences(video, nodes, connections).map((item) => item.kind).sort()).toEqual(["image", "text"]);
    });

    test("只有 storageKey 的图片节点也会作为预引用", () => {
        const video = videoNode("video");
        const image = imageNode("still", { storageKey: "resource:still" });
        const nodes = [video, image];
        const connections = [connection("still", "video")];
        expect(buildNodeMentionReferences(video, nodes, connections).map((item) => item.nodeId)).toEqual(["still"]);
    });

    test("空图片节点连入视频后仍显示预引用", () => {
        const video = videoNode("video");
        const image = imageNode("still", {});
        const nodes = [video, image];
        const connections = [connection("still", "video")];
        expect(buildNodeMentionReferences(video, nodes, connections).map((item) => ({ kind: item.kind, label: item.label }))).toEqual([{ kind: "image", label: "still" }]);
    });

    test("视频自己的连入素材不会被下游配置节点输入覆盖", () => {
        const video = videoNode("video");
        const image = imageNode("still");
        const text = textNode("note");
        const config = configNode("config");
        const nodes = [video, image, text, config];
        const connections = [connection("still", "video"), connection("note", "config"), connection("video", "config")];
        expect(buildNodeMentionReferences(video, nodes, connections).map((item) => item.nodeId).sort()).toEqual(["note", "still"]);
    });
});

describe("collectUpstreamVideoNodes", () => {
    test("下游视频节点能回溯到上游视频源", () => {
        const source = videoNode("source-video");
        const segment = videoNode("segment-video");
        const target = videoNode("target-video");
        const text = textNode("script");
        const nodes = [target, segment, source, text];
        const connections = [connection("source-video", "segment-video"), connection("segment-video", "target-video"), connection("script", "segment-video")];
        expect(collectUpstreamVideoNodes("target-video", nodes, connections).map((node) => node.id)).toEqual(["target-video", "segment-video", "source-video"]);
    });

    test("存在环时不会死循环", () => {
        const a = videoNode("a");
        const b = videoNode("b");
        const nodes = [a, b];
        const connections = [connection("a", "b"), connection("b", "a")];
        expect(collectUpstreamVideoNodes("a", nodes, connections).length).toBe(2);
    });
});
