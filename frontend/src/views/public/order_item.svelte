<script lang="ts">
  import { UserOrder } from "../../client/order";
  import { slide } from "svelte/transition";
  import { quintOut } from "svelte/easing";
  import { FormatPriceDefault } from "../../utils/format_price";
  import { onMount } from "svelte";

  export let order: UserOrder;
  let showOrderItems = false;
  let orderStatus = "";
  let orderDate = "";
  onMount(() => {
    let createdAt = new Date(order.createdAt);
    orderDate = createdAt.toISOString().slice(0, 10);
    orderDate += " " + createdAt.getHours() + ":" + createdAt.getMinutes();
    switch (order.status) {
      case "PROCESS":
        orderStatus = "в процессе";
        break;
      case "PAID":
        orderStatus = "оплачен";
        break;
      case "CANCELED":
        orderStatus = "отменён";
        break;
      case "SUCCESS":
        orderStatus = "выполнен";
        break;
      default:
        orderStatus = "";
        break;
    }
  });
</script>

<div class="order-box">
  <section class="order-info">
    <div class="order-text">
      Заказ от {orderDate}
    </div>
    <horizontalSpacer class="primary-bg small" />
    <button
      class="show-hide-order-items-button"
      on:click={() => {
        showOrderItems = !showOrderItems;
      }}>v</button
    >
  </section>

  <section>
    {#if showOrderItems}
      <div
        class="order-items primary-bg"
        transition:slide={{
          delay: 250,
          duration: 300,
          easing: quintOut,
          axis: "y",
        }}
      >
        {#each order.items as item}
          <div class="user-order-item">
            <div class="name-count">{item.name} x {item.count}</div>
            <horizontalSpacer class="primary-bg small" />
            <div class="total-price">
              {FormatPriceDefault(item.totalPrice)}
            </div>
          </div>
        {/each}
        <horizontalSpacer class="primary-bg small" />
        <div class="order-total">
          {#if orderStatus != ""}
            <div>Статус:</div>
            <horizontalSpacer class="primary-bg small" />
            <div class="order-text">{orderStatus}</div>
            <horizontalSpacer class="primary-bg small" />
          {/if}
          <div class="user-order-total">Итого:</div>
          <horizontalSpacer class="primary-bg small" />
          <div>{FormatPriceDefault(order.total)}</div>
        </div>
      </div>
    {/if}
  </section>
</div>

<style>
  .order-box {
    padding: 10px;
    --order-item-border-color: var(--primary-bg-color) * 0.1;
  }
  .order-total {
    display: flex;
    flex-direction: row;
  }
  .order-items {
    border-radius: 8px;
    padding: 10px;
    margin: 5px;
    border: 2px solid var(--primary-text-color);
    background-clip: padding-box;
  }
  .order-info {
    padding-top: 3px;
    display: flex;
    flex-direction: row;
    align-items: baseline;
    width: 100%;
    text-align: start;
    font-weight: normal;
  }
  .show-hide-order-items-button {
    display: inline-flex;
    justify-content: center;
    background-color: transparent;
    border: 1px solid var(--tg-theme-text-color);
    padding-top: 0;
    padding-bottom: 0;
    font-weight: normal;
  }
  .show-hide-order-items-button:hover {
    color: green;
  }
  .user-order-item {
    display: inline-flex;
    flex-direction: row;
    flex-wrap: nowrap;

    width: 100%;
  }
  .user-order-item .name-count {
    text-align: left;
  }
  .user-order-item .total-price {
    text-align: right;
  }
</style>
