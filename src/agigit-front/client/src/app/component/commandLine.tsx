import { Editor } from "@monaco-editor/react";
import { Button } from "@nextui-org/react";
import { useEffect, useState } from "react";
import { greenBtnCss } from "../constant/tailwind";
import { getLocalStorage, setLocalStorage } from "../util/localStorage";
import { registerAgiGit } from "../util/registerAgiGit";

interface IProps {
  setCommand: any;
}

export default function CommandLine(props: IProps) {
  const { setCommand } = props;
  const [editorValue, setEditorValue] = useState<string | undefined>("");
  const [recordList, setRecordList] = useState<any>([]);

  const handleEditorChange = (value: string | undefined) => {
    setEditorValue(value);
  };

  useEffect(() => {
    const list = getLocalStorage("record") ?? [];
    setRecordList(list);
  }, []);

  const codeArr = [
    "agigit relay add",
    "agigit relay remove",
    "agigit pull",
    "agigit push",
  ];

  return (
    <div className="flex y-10 px-4 my-4">
      <div className="w-2/3 rounded-xl border p-2 mr-2">
        <div className="">
          <label htmlFor="comment" className="sr-only">
            Add your code
          </label>
          <Editor
            height="40vh"
            defaultLanguage="AgiGit"
            value={editorValue}
            theme="AgiGitTheme"
            onMount={registerAgiGit}
            onChange={handleEditorChange}
          />
        </div>
        <div className="flex flex-row-reverse mt-2">
          <Button
            className={greenBtnCss}
            onClick={() => {
              if ((editorValue ?? "").length > 0) {
                setCommand(editorValue);
                const list = [editorValue, ...recordList];
                setRecordList(list);
                setEditorValue("");
                setLocalStorage("record", list);
              }
            }}
          >
            Run
          </Button>
        </div>
      </div>
      <div className="w-1/3 rounded-xl border p-4 bg-[#fff] font-bold overflow-y-scroll h-[400px] scrollbar-thin">
        <div className="text-sm text-bold mb-2">
          We currently recommend following instructions 😋
        </div>
        <div className="mb-2">
          {codeArr.map((item, idx) => {
            return (
              <div key={item}>
                <code className="text-sm">{`${idx+1}. ${item}`}</code>
              </div>
            );
          })}
        </div>
        <div className="text-xl text-bold mb-2">History Record</div>
        <div>
          {recordList.map((item: string, idx: number) => {
            return (
              <div
                className="text-sm w-full rounded-xl bg-[#F8F6E3] mb-2 p-2 font-normal"
                key={idx}
              >
                {idx + 1}: {item}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
