import {keyCodes} from 'keycodes';

/**
 * LayoutMaker
 */
export class LayoutMaker {
  constructor() {
    this.selectedKey = null;
    this.selectedLayer = null;
    this.layout = [];
    this.nrOfLayers = 0;
  }

  /**
   * addLayer
   * adds a layer to the layout
   */
  addLayer() {
    var l = this.nrOfLayers;
    this.layout[l] = [];
    var $template = $('.layer-template').clone(false);
    $template.attr('class', 'layer layer-'+l);
    $template.attr('data-layer', l);
    $('body').append($template);

    // Setup handlers of all the keys
    for(var i=0; i <= 80; i++) {
      d3.select('.layer.layer-'+l+' .key.key-'+i).on("mouseover", function() {
        d3.select(this).classed({highlight: true});
      });
      d3.select('.layer.layer-'+l+' .key.key-'+i).on('mouseout', function() {
        d3.select(this).classed({highlight: false});
      });
      d3.select('.layer.layer-'+l+' .key.key-'+i).on('click', this.selectKey.bind(this));
    }
    this.nrOfLayers++;
  }

  /**
   * starts the maker, adds the initial layer
   */
  start() {
    this.addLayer();

    d3.select('body').on('keydown', this.pressedKey.bind(this));
    d3.select('body').on('keyup', this.pressedKey.bind(this));
    d3.select('#save').on('click', this.save.bind(this));
    d3.select('#add-layer').on('click', this.addLayer.bind(this));
  }

  /**
   * Set key in layout
   * @param layer {Number} layer of the key
   * @param key {Number} key position of the key (within the layer)
   * @param keyCode {Number} keycode
   */
  setKey(layer, key, keyCode) {
    this.layout[this.selectedLayer][this.selectedKey-1] = keyCode;
  }

  /*
   * EVENT HANDLERS
   */

  /**
   * the user selects a key (using mouse)
   */
  selectKey() {
    var $this = $(d3.event.target);
    var key = $this.data('key');
    var layer = $this.closest('svg').data('layer');

    if (this.selectedLayer != null && this.selectedKey != null) {
      d3.select('.layer.layer-'+this.selectedLayer+' .key.key-'+this.selectedKey).classed({selected: false});
    }
    if (this.selectedKey != key || this.selectedLayer != layer) {
      d3.select(d3.event.target).classed({selected: true});
      this.selectedKey = key;
      this.selectedLayer = layer;
    } else {
      this.selectedKey = null;
      this.selectedLayer = null;
    }
  }

  /**
   * user pressed a key
   */
  pressedKey() {
    if(this.selectedKey != null && this.selectedLayer != null) {
      var $key = d3.select('.layer.layer-'+this.selectedLayer+' .key.key-'+this.selectedKey);
      var $text = d3.select('.layer.layer-'+this.selectedLayer+' .label.label-'+this.selectedKey);
      var $wrapper = $key.node().parentNode;
      if (!keyCodes[d3.event.keyCode]) {
        console.log("Key not recognised, please report.");
        console.log(d3.event);
        return;
      }

      if ($text.empty()) {
        $text = d3.select($wrapper).append('text')
          .attr('class', 'label label-'+this.selectedKey)
          .attr('x', +$key.attr('x') + $key.attr('width')/2)
          .attr('y', +$key.attr('y'));

        $text.append('tspan').attr('dx', 0).attr('dy', 30).html(keyCodes[d3.event.keyCode][0]);
      } else {
        $text.select('text tspan').html(keyCodes[d3.event.keyCode][0]);
      }
      this.setKey(this.selectedLayer, this.selectedKey, keyCodes[d3.event.keyCode][1]);
    }
    d3.event.preventDefault();
    return false;
  }

  /**
   * the user wants to save the layout
   */
  save() {
    // Each layout has 84 keys (14 rows 6 columns)
    var jsn = {
      "keyboard_layout": {
        "description": "dvorak",
        "layers": []
      }
    };

    for(var i=0; i<this.layout.length; i++) {
      var keymap = Array.apply(null, Array(84)).map(function (x, i) { return "KC_TRANSPARENT"; });
      for(var k=0; k< this.layout[i].length; k++) {
        keymap[k] = this.layout[i][k] || "KC_TRANSPARENT";
      }
      jsn['keyboard_layout']['layers'].push({"description": "", "keymap": keymap});

    }
    console.log(JSON.stringify(jsn, null, "  "));
  }
}
