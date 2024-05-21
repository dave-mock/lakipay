import java.io.FileInputStream;
import java.security.KeyStore;
import java.security.PrivateKey;

/**
 * tst
 */
public class tst {

    public static PrivateKey getPrivateKey() throws Exception {
        PrivateKey privateKey = null;
        try {
            String keyStorePath = "./IntellectARXSJWT.jks";
            String keyPassword = "Intellect01";
            String keyAlias = "intellect";
            KeyStore ks = KeyStore.getInstance(KeyStore.getDefaultType());
            ks.load(new FileInputStream(keyStorePath), keyPassword.toCharArray());
            KeyStore.PrivateKeyEntry keyEntry = (KeyStore.PrivateKeyEntry) ks.getEntry(keyAlias,
                    new KeyStore.PasswordProtection(keyPassword.toCharArray())); // alias name
            // from db or properties file.
            privateKey = keyEntry.getPrivateKey();
        } catch (Exception exp) {
            throw new Exception("Exception while reading the private key", exp);
        }
        return privateKey;
    }

    public static void main(String[] args) {
        System.out.println("Starting");
        try {
            PrivateKey key = getPrivateKey();
            System.out.println(key.getAlgorithm());
            System.out.println(key.getFormat());
            System.out.println(key.toString());
        } catch (Exception e) {
            // TODO Auto-generated catch block
            e.printStackTrace();
        }
        System.out.println();
    }
}