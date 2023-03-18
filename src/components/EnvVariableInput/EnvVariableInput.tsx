import PortInput from "../Input/PortInput";
import Input from "../Input/Input";
import { Vertical } from "../Layouts/Layouts";
import { EnvVariable } from "../../backend/backend";

type Props = {
    env: EnvVariable;
    value: any;
    onChange: (value: any) => void;
};

export default function EnvVariableInput(props: Props) {
    const { env, value, onChange } = props;

    const inputProps = {
        value,
        label: env.display_name,
        name: env.name,
        description: env.description,
        onChange: (e) => onChange(e.target.value),
    };

    let input;
    if (env.type === "port") {
        input = <PortInput {...inputProps} />;
    } else {
        input = <Input {...inputProps} />;
    }

    return <Vertical gap={6}>{input}</Vertical>;
}
