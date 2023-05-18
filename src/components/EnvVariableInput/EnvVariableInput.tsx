import PortInput from "../Input/PortInput";
import Input from "../Input/Input";
import { Vertical } from "../Layouts/Layouts";
import { EnvVariable } from "../../backend/backend";

type Props = {
    env: EnvVariable;
    value: any;
    onChange: (value: any) => void;
    disabled?: boolean;
};

export default function EnvVariableInput(props: Props) {
    const { env, value, onChange, disabled } = props;

    const inputProps = {
        value,
        label: env.display_name,
        name: env.name,
        description: env.description,
        onChange: (e) => onChange(e.target.value),
        disabled,
    };

    let input;
    if (env.type === "port") {
        input = <PortInput {...inputProps} />;
    } else {
        input = <Input {...inputProps} />;
    }

    return <Vertical gap={6}>{input}</Vertical>;
}