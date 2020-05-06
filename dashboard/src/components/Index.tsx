import * as React from "react";

import IndexForm from "./forms/IndexForm";
import SearchForm from "./forms/SearchForm";
import Panels from "./panels/Panels";
import Panes from "./panes/Panes";
import IndexTable from "./tables/IndexTable";
import Tabs from "./tabs/Tabs";
import IndexView from "./views/IndexView";
import ConsoleView from "./views/ConsoleView";

import logger from "../lib/logger";

function Index(props) {
    logger.info("Index.render", props.Kind, props);
    switch (props.Kind) {
        case "Msg":
            return <div>Msg</div>;
        case "Panels":
            return <Panels index={props} />;
        case "Tabs":
            return <Tabs index={props} />;
        case "Panes":
            return <Panes index={props} />;
        case "Table":
            return <IndexTable index={props} />;
        case "View":
            return <IndexView index={props} />;
        case "ConsoleView":
            return <ConsoleView index={props} />;
        case "SearchForm":
            return <SearchForm index={props} />;
        case "Form":
            return <IndexForm index={props} />;
        default:
            return <div>Unsupported Kind: {props.Kind}</div>;
    }
}

export default Index;
