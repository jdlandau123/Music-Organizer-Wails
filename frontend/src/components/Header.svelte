<script>
  import icon from "../assets/icons/settings_white_36dp.svg";
  import { SelectDirectory, SetConfig, GetConfig } from '../../wailsjs/go/main/App.js'
  import { onMount } from 'svelte';

  let showMenu = false;
  let saved = false;

  let config = {
    CollectionPath: '',
    DevicePath: ''
  };

  function setCollectionPath() {
    SelectDirectory().then(path => config.CollectionPath = path);
  }

  function setDevicePath() {
    SelectDirectory().then(path => config.DevicePath = path);
  }

  function setConfig() {
    SetConfig(config).then(() => {
      saved = true;
      getConfig();
      setTimeout(() => {
        saved = false;
        showMenu = false;
      }, 3000);
    });
  }

  function getConfig() {
    GetConfig().then(c => config = c);
  }

  onMount(() => getConfig());
</script>

<div class="wrapper">
  <h2>Music Organizer</h2>
  <button class="iconBtn" on:click={() => showMenu = !showMenu}>
    <img src={icon} alt="settings" />
  </button>
</div>
{#if showMenu}
  <div class="configMenu">
    <h3>Settings</h3>
    <div class="configSetting">
      Collection:
      <button on:click={setCollectionPath} class="configBtn">
        {config.CollectionPath ? config.CollectionPath : 'Select Collection Directory'}
      </button>
    </div>
    <div class="configSetting">
      Device:
      <button on:click={setDevicePath} class="configBtn">
        {config.DevicePath ? config.DevicePath : 'Select Device Directory'}
      </button>
    </div>
    <button class="primaryBtn" on:click={setConfig}>Save Changes</button>
    {#if saved}
      <p style="color: #4ba346;">Settings Saved!</p>
    {/if}
  </div>
{/if}

<style>
  .wrapper {
    display: flex;
    justify-content: space-between;
    padding: 10px;
    background-color: #4ba346;
    color: white;
    width: 100%;
  }

  .iconBtn {
    background-color: transparent;
    border: transparent;
    text-decoration: none;
    margin-right: 10px;
  }

  .iconBtn:hover {
    cursor: pointer;
  }

  .configMenu {
    position: absolute;
    top: 100px;
    right: 10px;
    z-index: 1;
    width: 30%;
    border: 1pt solid black;
    border-radius: 8px;
    background-color: white;
    padding: 10px;
  }

  h3 {
    font-weight: 500;
  }

  .configBtn {
    border: transparent;
    background-color: #fff;
    border-radius: 8px;
    padding: 10px;
  }

  .configSetting {
    display: flex;
    gap: 20px;
    align-items: center;
    margin-bottom: 10px;
  }

  .configSetting > *:hover {
    cursor: pointer;
    background-color: lightgrey;
  }
</style>
