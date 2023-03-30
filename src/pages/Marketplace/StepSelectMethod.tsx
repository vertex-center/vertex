import styles from "./Marketplace.module.sass";
import { Title } from "../../components/Text/Text";
import Button from "../../components/Button/Button";
import { DownloadMethod } from "./Marketplace";
import { SiGit } from "@icons-pack/react-simple-icons";

type StepSelectMethodProps = {
    method: DownloadMethod;
    onMethodChange: (method: DownloadMethod) => void;
    onNextStep: () => void;
};

export default function StepSelectMethod(props: StepSelectMethodProps) {
    const { method, onMethodChange, onNextStep } = props;

    return (
        <div className={styles.step}>
            <div className={styles.stepTitle}>
                <Title>Installation method</Title>
            </div>
            <div className={styles.buttons}>
                <Button
                    className={styles.button}
                    onClick={() => onMethodChange("marketplace")}
                    leftSymbol="precision_manufacturing"
                    selectable
                    selected={method === "marketplace"}
                >
                    <div className={styles.buttonContent}>
                        <div>Marketplace</div>
                        <div className={styles.buttonDescription}>
                            Download services from our online and certified
                            repository.
                        </div>
                    </div>
                </Button>
                <Button
                    className={styles.button}
                    onClick={() => onMethodChange("git")}
                    leftSymbol={<SiGit />}
                    selectable
                    selected={method === "git"}
                >
                    <div className={styles.buttonContent}>
                        <div>Git</div>
                        <div className={styles.buttonDescription}>
                            Clone services from GitHub, GitLab...
                        </div>
                    </div>
                </Button>
                <Button
                    className={styles.button}
                    onClick={() => onMethodChange("localstorage")}
                    leftSymbol="storage"
                    selectable
                    selected={method === "localstorage"}
                >
                    <div className={styles.buttonContent}>
                        <div>Local storage</div>
                        <div className={styles.buttonDescription}>
                            Point to vertex the path of an existing service on
                            your computer. Vertex will keep it installed there.
                        </div>
                    </div>
                </Button>
            </div>
            <Button
                primary
                large
                disabled={method === undefined}
                rightSymbol="navigate_next"
                onClick={onNextStep}
            >
                Next
            </Button>
        </div>
    );
}
