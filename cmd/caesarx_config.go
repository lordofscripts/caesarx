/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cmd

import (
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/cmn/prefs"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// Organization name. The application root config dir.
	ORGANIZATION string = "coralys"
	// This application name. A subdirectory of the organization name
	APPLICATION string = "caesarx"
	// Base name of the configuration file in ~/<user_config>/ORG/APP/
	CONFIG_BASE_FILENAME string = "caesarx.yaml"
	// Configuration schema version (for compatibility)
	CONFIG_SCHEMA_VERSION string = "1.0"

	DEFAULT_UNSET_NGRAM int = 0
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// configuration handler. Make sure to call InitConfiguration()
// after the type is instantiated.
type CaesarxConfig struct {
	Configuration *Config
	isGood        bool
}

// configuration model
type Config struct {
	Version  string             `yaml:"version"`
	Defaults *ConfigDefaults    `yaml:"defaults"`
	Profiles []*prefs.Recipient `yaml:"profiles,omitempty"`
}

type ConfigDefaults struct {
	AlphaName  string `yaml:"alphabet"`
	SlaveName  string `yaml:"supplementary"`
	NGramSize  int    `yaml:"ngram_size"`
	CipherName string `yaml:"preferred_cipher"`
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// instantiates a new configuration handler with a default.
func NewConfiguration() *CaesarxConfig {
	return &CaesarxConfig{
		Configuration: newDefaultConfig(),
		isGood:        false,
	}
}

func newDefaultConfig() *Config {
	var SampleProfiles []*prefs.Recipient = []*prefs.Recipient{
		prefs.NewProfileWithCipher("you+c@bitbucket.com", "Sample profile 1", caesarx.CaesarCipher, cmn.ISO_EN, cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED, &prefs.CaesarModel{Key: 'C'}),
		prefs.NewProfileWithCipher("you+d@bitbucket.com", "Sample profile 2", caesarx.DidimusCipher, cmn.ISO_ES, cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED, &prefs.CaesarModel{Key: 'D', Offset: 8}),
		prefs.NewProfileWithCipher("you+f@bitbucket.com", "Sample profile 3", caesarx.FibonacciCipher, cmn.ISO_PT, cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED, &prefs.CaesarModel{Key: 'F', Offset: 3}),
		prefs.NewProfileWithCipher("you+b@bitbucket.com", "Sample profile 4", caesarx.BellasoCipher, cmn.ISO_DE, cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED, &prefs.SecretsModel{Secret: "Ein Geheim"}),
		prefs.NewProfileWithCipher("you+v@bitbucket.com", "Sample profile 5", caesarx.VigenereCipher, cmn.ISO_IT, cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED, &prefs.SecretsModel{Secret: "Buongiorno a tutti"}),
		prefs.NewProfileWithCipher("you+a@bitbucket.com", "Sample profile 6", caesarx.AffineCipher, cmn.ISO_GR, cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED, &prefs.AffineModel{A: 7, B: 12, Ap: 21}),
	}

	return &Config{
		Version: CONFIG_SCHEMA_VERSION,
		Defaults: &ConfigDefaults{
			AlphaName:  cmn.ALPHA_NAME_ENGLISH,
			SlaveName:  string(OPT_SLAVE_NONE),
			NGramSize:  DEFAULT_UNSET_NGRAM,
			CipherName: caesarx.CaesarCipher.String(),
		},
		Profiles: SampleProfiles,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Initialize the user configuration. Read it an existing one or create.
func (c *CaesarxConfig) InitConfiguration() error {
	var err error = nil
	// get the platform-dependent configuration directory for this app
	cfgDir := app.GetConfigDir(ORGANIZATION, APPLICATION)
	// fully-qualified configuration filename
	cfgFile := path.Join(cfgDir, CONFIG_BASE_FILENAME)
	// do we already have one in the file system?
	if err = app.CheckFileExistsAndReadable(cfgFile); err != nil {
		// no, try to ensure we can access the configuration path
		if err = app.EnsureConfigDir(cfgDir); err == nil {
			// create a default configuration as fallback
			// attempt to save it to the user configuration
			err = c.SaveConfig()
		}
	} else {
		// read the existing configuration file
		err = c.readConfig(cfgFile)
	}

	return err
}

// whether we have a succesfully read user configuration.
// it is false if it was created on the spot.
func (c *CaesarxConfig) IsGood() bool {
	return c.isGood
}

// reads the user configuration file
func (c *CaesarxConfig) readConfig(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		mlog.ErrorT("read-config", mlog.At(), mlog.Err(err))
		return err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		mlog.ErrorT("unmarshall-config", mlog.At(), mlog.Err(err))
		return err
	}

	c.Configuration = &config
	c.isGood = true
	return nil
}

// save the user configuration to a file
func (c *CaesarxConfig) SaveConfig() error {
	data, err := yaml.Marshal(&c.Configuration)
	if err != nil {
		mlog.ErrorT("marshall-config", mlog.At(), mlog.Err(err))
		return err
	}

	cfgDir := app.GetConfigDir(ORGANIZATION, APPLICATION)
	cfgFile := path.Join(cfgDir, CONFIG_BASE_FILENAME)

	err = os.WriteFile(cfgFile, data, 0644)
	if err != nil {
		mlog.ErrorT("write-config", mlog.At(), mlog.Err(err))
		return err
	}

	return nil
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
