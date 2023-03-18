import { Fragment, useEffect, useState } from "react";
import { Title } from "../../components/Text/Text";
import { getInstance, Instance } from "../../backend/backend";
import { useParams } from "react-router-dom";
import EnvVariableInput from "../../components/EnvVariableInput/EnvVariableInput";
import Button from "../../components/Button/Button";

type Props = {};

export default function BayDetailsEnv(props: Props) {
    const { uuid } = useParams();

    const [env, setEnv] = useState<any[]>();

    const [instance, setInstance] = useState<Instance>();

    useEffect(() => {
        setEnv(
            instance?.environment.map((e) => ({
                env: e,
                value: instance?.env[e.name] ?? e.default ?? "",
            }))
        );
    }, [instance?.environment]);

    const onChange = (i: number, value: any) => {
        setEnv((prev) =>
            prev.map((el, index) => {
                if (index !== i) return el;
                return { ...el, value };
            })
        );
    };

    useEffect(() => {
        getInstance(uuid).then((i: Instance) => setInstance(i));
    }, [uuid]);

    const save = () => {};

    return (
        <Fragment>
            <Title>Environment</Title>
            {env?.map((env, i) => (
                <EnvVariableInput
                    env={env.env}
                    value={env.value}
                    onChange={(v) => onChange(i, v)}
                />
            ))}
            <Button primary large onClick={save} rightSymbol="save">
                Save
            </Button>
        </Fragment>
    );
}
