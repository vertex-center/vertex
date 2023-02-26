import { HTMLProps } from "react";

type Props = HTMLProps<HTMLInputElement> & {};

export default function Input(props: Props) {
    const { ...others } = props;

    return <input {...others} />;
}
