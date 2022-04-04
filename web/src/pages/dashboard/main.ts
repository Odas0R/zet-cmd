import { TW } from "@utils/tailwind-mixin";
import type { TemplateResult } from "lit";
import { html, LitElement } from "lit";
import { customElement } from "lit/decorators.js";

import "@components/button-counter";
import "@styles/main.css";

@customElement("index-page")
export class IndexPage extends TW(LitElement) {
  render(): TemplateResult {
    return html`
      <div class="container">
        <button-counter></button-counter>
      </div>
    `;
  }
}
