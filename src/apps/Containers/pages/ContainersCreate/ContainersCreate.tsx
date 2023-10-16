import styles from "./ContainersCreate.module.sass";
import { BigTitle, Title } from "../../../../components/Text/Text";
import Input from "../../../../components/Input/Input";
import { Horizontal, Vertical } from "../../../../components/Layouts/Layouts";
import NoItems from "../../../../components/NoItems/NoItems";
import Button from "../../../../components/Button/Button";
import classNames from "classnames";
import Spacer from "../../../../components/Spacer/Spacer";

type Props = {};

export default function ContainersCreate(props: Readonly<Props>) {
    return (
        <Vertical gap={30}>
            <BigTitle className={styles.bigTitle}>Create container</BigTitle>
            <Vertical gap={25}>
                <Title className={styles.title}>Info</Title>
                <div className={classNames(styles.content, styles.inputs)}>
                    <Input label="Service ID" />
                    <Input label="Service name" />
                    <Input label="Repository" />
                    <Input label="Description" />
                    <Input label="Color" />
                </div>
            </Vertical>

            <Vertical gap={10}>
                <Horizontal className={styles.title} alignItems="center">
                    <Title>Environment</Title>
                    <Spacer />
                    <Button primary rightIcon="add">
                        Add ENV variable
                    </Button>
                </Horizontal>
                <Vertical className={styles.content} gap={10}>
                    <NoItems icon="list" text="No environment variables" />
                    <div></div>
                </Vertical>
            </Vertical>

            <Vertical gap={10}>
                <Horizontal className={styles.title} alignItems="center">
                    <Title>URLs</Title>
                    <Spacer />
                    <Button primary rightIcon="add">
                        Add URL
                    </Button>
                </Horizontal>
                <Vertical className={styles.content} gap={10}>
                    <NoItems icon="public" text="No URLs." />
                </Vertical>
            </Vertical>

            <Vertical gap={25}>
                <Title className={styles.title}>Docker installation</Title>
                <Vertical
                    className={classNames(styles.content, styles.inputs)}
                    gap={15}
                >
                    <Input label="Docker image" required />
                    <Input label="Command" />
                    <Input label="Ports" />
                    <Input label="Volumes" />
                    <Input label="Environment" />
                    <Input label="Capabilities" />
                    <Input label="Sysctls" />
                </Vertical>
            </Vertical>
        </Vertical>
    );
}
