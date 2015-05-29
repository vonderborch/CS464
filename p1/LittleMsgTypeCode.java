
/*
  WARNING: THIS FILE IS AUTO-GENERATED. DO NOT MODIFY.

  This file was generated from .idl using "rtiddsgen".
  The rtiddsgen tool is part of the RTI Connext distribution.
  For more information, type 'rtiddsgen -help' at a command shell
  or consult the RTI Connext manual.
*/
    
import com.rti.dds.typecode.*;


public class LittleMsgTypeCode {
    public static final TypeCode VALUE = getTypeCode();

    private static TypeCode getTypeCode() {
        TypeCode tc = null;
        int i=0;
        StructMember sm[] = new StructMember[3];

        sm[i]=new StructMember("messagenumber",false,(short)-1,false,(TypeCode)TypeCode.TC_LONG,0,false); i++;
        sm[i]=new StructMember("sender",false,(short)-1,false,(TypeCode)new TypeCode(TCKind.TK_STRING,(MSG_LEN.VALUE)),1,false); i++;
        sm[i]=new StructMember("message",false,(short)-1,false,(TypeCode)new TypeCode(TCKind.TK_STRING,(MSG_LEN.VALUE)),2,false); i++;

        tc = TypeCodeFactory.TheTypeCodeFactory.create_struct_tc("LittleMsg",ExtensibilityKind.EXTENSIBLE_EXTENSIBILITY,sm);
        return tc;
    }
}
