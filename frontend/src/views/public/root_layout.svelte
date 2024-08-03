<script>
  import { navigateTo } from "svelte-router-spa";
  import { UserIsAdmin } from "../../client/user";
  import { onMount } from "svelte";

  onMount(async () => {
    let isAdmin = await UserIsAdmin();
    if (isAdmin) {
      navigateTo("/admin");
    } else {
      navigateTo("/dishes");
    }
  });
</script>

<div></div>
