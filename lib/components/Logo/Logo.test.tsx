import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import { Logo } from "./Logo";

test("it renders", () => {
    render(<Logo />);
    const logo = screen.getByAltText("App Logo");
    expect(logo).toBeInTheDocument();
});
