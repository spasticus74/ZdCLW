<script>
  import Tabs from "./Tabs.svelte";
  import ContactsTab from "./ContactsTab.svelte";
  import AddCallTab from "./AddCallTab.svelte";
  import LoginTab from "./LoginTab.svelte";

  // List of tab items with labels, values and assigned components
  let items = [
    { label: "Add Call",
      value: 1,
      component: AddCallTab
    },
    { label: "Contacts",
      value: 2,
      component: ContactsTab
    },
    { label: "Login",
      value: 3,
      component: LoginTab
    }
  ];
</script>

<Tabs {items} />