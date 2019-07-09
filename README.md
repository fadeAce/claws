# claws
### made for communicating with mainstream chain in one interface
<br>

#### core conception
support a chain needs coding
<br>support a type of chain based crypto needs coding
<br>support a chain based crypto needs not more than configuration

#### core usage

    // step.1 load config file
    conf := &types.Claws{}
    cfg, _ := ioutil.ReadFile("./claws.yml")
    _ = yaml.Unmarshal(cfg, conf)
    
    // step.2 setup builder with config file
    SetupGate(conf, nil)

    // step.3 just use wallet from buildered one
    wallet := claws.Builder.BuildWallet("eth")
    
    // step.4 do sth with wallet , like fetching txs from a block
    txns, err := wallet.UnfoldTxs(conf.Ctx, big.NewInt(4356126))
    