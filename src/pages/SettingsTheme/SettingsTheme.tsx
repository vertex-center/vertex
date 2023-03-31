import { Fragment, useContext } from "react";
import Button from "../../components/Button/Button";
import { Vertical } from "../../components/Layouts/Layouts";
import { Title } from "../../components/Text/Text";
import { ThemeContext } from "../../index";

type Props = {};

export default function SettingsTheme(props: Props) {
    const { theme, setTheme } = useContext(ThemeContext);

    return (
        <Fragment>
            <Title>Theme</Title>
            <Vertical gap={8}>
                <Button
                    onClick={() => setTheme("vertex-dark")}
                    selectable
                    selected={theme === "vertex-dark"}
                >
                    Vertex Dark
                </Button>
                <Button
                    onClick={() => setTheme("vertex-light")}
                    selectable
                    selected={theme === "vertex-light"}
                >
                    Vertex Light
                </Button>
            </Vertical>
        </Fragment>
    );
}
