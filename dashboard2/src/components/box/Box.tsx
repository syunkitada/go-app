import data from "../../data";
import converter from "../../lib/converter";
import locationData from "../../data/locationData";
import service from "../../apps/service";
import Table from "../table/Table";
import SearchForm from "../form/SearchForm";
import LineGraphCard from "../card/LineGraphCard";

export function Render(input: any) {
    const { id, View } = input;
    const keyPrefix = `${id}-Box-`;
    const location = locationData.getLocationData();
    const boxData = data.service.data[View.DataKey];

    const panelsId = `${keyPrefix}panels`;

    $(`#${id}`).html(`
    <div>
      <div id="${panelsId}"></div>
    </div>
    `);

    function getFieldHtml(input: any) {
        const { field, data } = input;
        if (!data) {
            return `<tr><td>${field.Name}</td><td></td>`;
        }

        let html: any;
        let subData: any;
        if (field.DataKey) {
            subData = data[field.DataKey];
        } else {
            subData = data[field.Name];
        }

        switch (field.Kind) {
            default:
                html = `${subData}`;
        }
        return `<tr><td>${field.Name}</td><td>${html}</td>`;
    }

    function renderPanels() {
        const renderHandlers = [];
        const panelsGroups = [];
        for (let i = 0, len = View.PanelsGroups.length; i < len; i++) {
            const panelsGroup = View.PanelsGroups[i];
            switch (panelsGroup.Kind) {
                case "Cards":
                    const cards = [];
                    for (
                        let j = 0, jlen = panelsGroup.Cards.length;
                        j < jlen;
                        j++
                    ) {
                        const card = panelsGroup.Cards[j];
                        let cardData = boxData;
                        if (card.SubDataKey) {
                            cardData = boxData[card.SubDataKey];
                        }
                        switch (card.Kind) {
                            case "Fields":
                                const fields = [];
                                for (
                                    let x = 0, xlen = card.Fields.length;
                                    x < xlen;
                                    x++
                                ) {
                                    const field = card.Fields[x];
                                    fields.push(
                                        getFieldHtml({ field, data: cardData })
                                    );
                                }
                                cards.push(`
                                    <div class="col m6">
                                      <h4>${card.Name}</h4>
                                      <table class="table">
                                        <thead><tr><th>Field Name</th><th>Field Value</th></tr></thead>
                                        <tbody>${fields.join("")}</tbody>
                                      </table>
                                    </div>
                                `);

                                break;
                            case "Table":
                                cards.push(`
                                    <div class="col m6">
                                      <h5>${card.Name}</h5>
                                      <div id="${keyPrefix}${card.Name}"></div>
                                    </div>
                                `);

                                let tableData: any = [];
                                if (cardData[card.DataKey]) {
                                    tableData = cardData[card.DataKey];
                                }

                                renderHandlers.push({
                                    render: Table.Render,
                                    input: {
                                        id: `${keyPrefix}${card.Name}`,
                                        View: card,
                                        tableData
                                    }
                                });
                                break;
                            default:
                                cards.push(
                                    `<div class="col m6">UnknownCard: ${card.Kind}</div>`
                                );
                                break;
                        }
                    }
                    panelsGroups.push(
                        `<div class="row">${cards.join("")}</div>`
                    );
                    break;

                case "SearchForm":
                    panelsGroups.push(
                        `<div class="row">
                            <h5>${panelsGroup.Name}</h5>
                            <div id="${keyPrefix}searchForm-${panelsGroup.Name}">
                            </div>
                        </div>`
                    );

                    renderHandlers.push({
                        render: SearchForm.Render,
                        input: {
                            id: `${keyPrefix}searchForm-${panelsGroup.Name}`,
                            View: panelsGroup,
                            onSubmit: function (input: any) {
                                const { searchQueries } = input;
                                const newLocation = Object.assign(
                                    {},
                                    location,
                                    {
                                        SearchQueries: searchQueries
                                    }
                                );
                                service.getQueries({
                                    view: { id, View },
                                    location: newLocation
                                });
                            }
                        }
                    });

                    break;

                case "MetricsGroups":
                    const metricsGroups = boxData[panelsGroup.DataKey];
                    if (!metricsGroups) {
                        continue;
                    }
                    for (
                        let j = 0, jlen = metricsGroups.length;
                        j < jlen;
                        j++
                    ) {
                        const metricsGroup = metricsGroups[j];
                        const cards: any = [];
                        if (!metricsGroup.Metrics) {
                            continue;
                        }
                        for (
                            let x = 0, xlen = metricsGroup.Metrics.length;
                            x < xlen;
                            x++
                        ) {
                            const metric = metricsGroup.Metrics[x];
                            cards.push(`
                                <div class="col m6">
                                <h5>${metric.Name}</h5>
                                <div id="${keyPrefix}metric-${converter.escape_id(
                                metric.Name
                            )}"></div></div>
                            `);
                            renderHandlers.push({
                                render: LineGraphCard.Render,
                                input: {
                                    id: `${keyPrefix}metric-${converter.escape_id(
                                        metric.Name
                                    )}`,
                                    metric: metric
                                }
                            });
                        }
                        panelsGroups.push(`
                          <div class="row">
                            ${cards.join("")}
                          </div>
                        `);
                    }
                    break;

                default:
                    panelsGroups.push(
                        `<div class="row"><div class="col m6">UnknownPanels: ${panelsGroup.Kind}</div></div>`
                    );
            }
        }
        $(`#${panelsId}`).html(panelsGroups.join(""));

        for (let i = 0, len = renderHandlers.length; i < len; i++) {
            const handler = renderHandlers[i];
            handler.render(handler.input);
        }
    }
    renderPanels();
}

const index = {
    Render
};
export default index;
