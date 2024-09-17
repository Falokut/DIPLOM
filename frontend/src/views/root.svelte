<script>
  import { navigate } from "svelte-routing";
  import { UserIsAdmin } from "../client/user";
  import { onMount } from "svelte";

  onMount(async () => {
    let isAdmin = await UserIsAdmin();
    if (isAdmin) {
      navigate("/admin");
    } else {
      navigate("/public-index");
    }
  });
</script>

<div></div>
