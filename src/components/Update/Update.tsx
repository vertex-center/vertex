import { PropsWithChildren, useState } from "react";

import styles from "./Update.module.sass";
import Button from "../Button/Button";
import Spacer from "../Spacer/Spacer";
import Symbol from "../Symbol/Symbol";
import { Horizontal } from "../Layouts/Layouts";
import Progress from "../Progress";

export function Updates(props: PropsWithChildren) {
    return <div {...props} />;
}

type Props = {
    name: string;
    onUpdate: () => Promise<void>;
    available?: boolean;

    current_version?: string;
    latest_version?: string;
};

export default function Update(props: Props) {
    const { name, available, current_version, latest_version } = props;

    const [isLoading, setIsLoading] = useState(false);

    const onUpdate = () => {
        setIsLoading(true);
        props.onUpdate().finally(() => setIsLoading(false));
    };

    return (
        <div className={styles.update}>
            <div className={styles.name}>{name}</div>
            <Spacer />
            {available ? (
                <Horizontal gap={20} alignItems="center">
                    {current_version && latest_version && (
                        <code>
                            {current_version} {"->"} {latest_version}
                        </code>
                    )}
                    {!isLoading && (
                        <Button rightSymbol="download" onClick={onUpdate}>
                            Install
                        </Button>
                    )}
                    {isLoading && <Progress infinite />}
                </Horizontal>
            ) : (
                <div className={styles.status}>
                    <Symbol name="check" />
                    Up-to-date
                </div>
            )}
        </div>
    );
}
